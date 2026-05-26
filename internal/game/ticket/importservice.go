package ticket

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

type TicketImportJob struct {
	ID          uint      `gorm:"primaryKey"`
	CampaignID  uint      `gorm:"not null"`
	Status      string    `gorm:"size:30;default:PENDING"`
	TotalRows   int       `gorm:"default:0"`
	SuccessRows int       `gorm:"default:0"`
	FailedRows  int       `gorm:"default:0"`
	ErrorLog    string    `gorm:"type:text"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TicketImportJob) TableName() string {
	return "ticketimportjobs"
}

type RowErrorDetail struct {
	Row   int    `json:"row"`
	Value string `json:"value"`
	Reason string `json:"reason"`
}

type ImportService interface {
	Import(ctx context.Context, req ImportExcelRequest) (*TicketImportJob, error)
	ProcessExcelAsync(ctx context.Context, job *TicketImportJob, filePath string)
	insertBatchWithMetrics(ctx context.Context, tickets []Ticket) int
	failJob(ctx context.Context, jobID uint, reason string)
}

type importService struct {
	ticketRepo Repository
}

func NewImportService(ticketRepo Repository) ImportService{
	return &importService{ticketRepo: ticketRepo}
}

func (s *importService) Import(ctx context.Context, req ImportExcelRequest) (*TicketImportJob, error) {
	src, err := req.File.Open()
	if err != nil {
		return nil, fmt.Errorf("không thể mở file nguồn từ request: %w", err)
	}
	defer src.Close()

	tempFile, err := os.CreateTemp("", "ticket-import-*.xlsx")
	if err != nil {
		return nil, fmt.Errorf("thất bại khi tạo không gian file tạm: %w", err)
	}
	tempFilePath := tempFile.Name()
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, src); err != nil {
		_ = os.Remove(tempFilePath)
		return nil, fmt.Errorf("thất bại khi ghi dữ liệu vào file tạm: %w", err)
	}

	job := &TicketImportJob{
		CampaignID: req.CampaignID,
		Status: "PENDING",
	}

	if err := s.ticketRepo.CreateJob(ctx, job); err != nil {
		_ = os.Remove(tempFilePath)
		return nil, fmt.Errorf("thất bại khi khởi tạo tiến độ tác vụ dưới DB: %w", err)
	}

	go s.ProcessExcelAsync(context.Background(), job, tempFilePath)

	return job, nil
}

func (s *importService) ProcessExcelAsync(ctx context.Context, job *TicketImportJob, filePath string) {	
	defer os.Remove(filePath)

	// 1. Chuyển trạng thái sang PROCESSING thông qua Repo
	if err := s.ticketRepo.UpdateJobStatus(ctx, job.ID, "PROCESSING"); err != nil {
		slog.Error("Không thể cập nhật trạng thái Job sang PROCESSING", "job_id", job.ID, "err", err)
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		slog.Error("Không mở được file excel", "job_id", job.ID, "err", err)
		s.failJob(ctx, job.ID, fmt.Sprintf("Không mở được file excel: %v", err))
		return
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.Rows(sheetName)
	if err != nil {
		s.failJob(ctx, job.ID,fmt.Sprintf("Không đọc được nội dung Sheet: %v", err))
		return
	}
	defer rows.Close()

	seenInFile := make(map[string]int)
	var errorDetails []RowErrorDetail
	var batchTickets []Ticket

	const batchSize = 2000
	rowCount := 0
	totalSuccess := 0
	totalFailed := 0

	for rows.Next() {
		rowCount++
		if rowCount == 1 {
			continue
		}

		row, err := rows.Columns()
		if err != nil || len(row) < 2 {
			totalFailed++
			errorDetails = append(errorDetails, RowErrorDetail{Row: rowCount, Value: "", Reason: "Dòng dữ liệu không đủ 2 cột (Mã vé và Username)"})
			continue
		}

		ticketNumber := row[0]
		username := row[1]

		if ticketNumber == "" {
			totalFailed++
			errorDetails = append(errorDetails, RowErrorDetail{Row: rowCount, Value: "", Reason: "Mã vé trống"})
			continue
		}
		if len(ticketNumber) > 21 {
			totalFailed++
			errorDetails = append(errorDetails, RowErrorDetail{Row: rowCount, Value: ticketNumber, Reason: "Mã vé vượt quá 21 ký tự quy định"})
			continue
		}

		if username == "" {
			totalFailed++
			errorDetails = append(errorDetails, RowErrorDetail{Row: rowCount, Value: "", Reason: "Username trống"})
			continue
		}
		if len(username) > 100 {
			totalFailed++
			errorDetails = append(errorDetails, RowErrorDetail{Row: rowCount, Value: username, Reason: "Username vượt quá 100 ký tự"})
			continue
		}

		if firstSeenRow, exists := seenInFile[ticketNumber]; exists {
			totalFailed++
			errorDetails = append(errorDetails, RowErrorDetail{
				Row: rowCount,
				Value: ticketNumber,
				Reason: fmt.Sprintf("Mã vé bị trùng lặp với dòng số %d trong chính file này", firstSeenRow),
			})
			continue
		}

		seenInFile[ticketNumber] = rowCount

		batchTickets = append(batchTickets, Ticket{
			CampaignID:   job.CampaignID,
        	TicketNumber: ticketNumber,
        	Username:     username,
        	IsWinner:     false,
        	IsCanceled:   false,
		})

		if len(batchTickets) >= batchSize {
			rowsInserted := s.insertBatchWithMetrics(ctx, batchTickets)
			totalSuccess += rowsInserted

			dbDuplicates := len(batchTickets) - rowsInserted
			if dbDuplicates > 0 {
				totalFailed += dbDuplicates

				errorDetails = append(errorDetails, RowErrorDetail{
					Row: rowCount,
					Value:  fmt.Sprintf("Mẻ dòng từ %d", rowCount-batchSize+1),
					Reason: fmt.Sprintf("Có %d mã vé bị trùng lặp với dữ liệu đã tồn tại trong hệ thống (DB)", dbDuplicates),
				})
			}

			batchTickets = batchTickets[:0]
			_ = s.ticketRepo.UpdateJobProgress(ctx, job.ID, rowCount-1, totalSuccess, totalFailed)
		}
	}

	if len(batchTickets) > 0 {
		rowsInserted := s.insertBatchWithMetrics(ctx, batchTickets)
		totalSuccess += rowsInserted
		dbDuplicates := len(batchTickets) - rowsInserted
		if dbDuplicates > 0 {
			totalFailed += dbDuplicates
			errorDetails = append(errorDetails, RowErrorDetail{
				Row:    rowCount,
				Value:  "Mẻ cuối",
				Reason: fmt.Sprintf("Có %d mã vé bị trùng lặp với dữ liệu đã tồn tại trong hệ thống (DB)", dbDuplicates),
			})
		}
	}

	errorLogJSON, _ := json.Marshal(errorDetails)
	if err := s.ticketRepo.CompleteJob(ctx, job.ID, rowCount-1, totalSuccess, totalFailed, string(errorLogJSON)); err != nil {
		slog.Error("Thất bại khi cập nhật trạng thái kết thúc Job", "job_id", job.ID, "err", err)
	}
}

func (s *importService) insertBatchWithMetrics(ctx context.Context, tickets []Ticket) int {
	if len(tickets) == 0 {
		return 0
	}
	rowsAffected, err := s.ticketRepo.CreateInBatchesWithCount(ctx, tickets)
	if err != nil {
		slog.Error("Thất bại nặng khi chèn mẻ dữ liệu", "err", err)
		return 0
	}
	return rowsAffected
}

func (s *importService) failJob(ctx context.Context, jobID uint, reason string) {
	errorLog := fmt.Sprintf(`[{"reason": "%s"}]`, reason)
	if err := s.ticketRepo.CompleteJob(ctx, jobID, 0, 0, 0, errorLog); err != nil {
		slog.Error("Thất bại khi cập nhật trạng thái FAILED cho Job", "job_id", jobID, "err", err)
	}
	_ = s.ticketRepo.UpdateJobStatus(ctx, jobID, "FAILED")
}