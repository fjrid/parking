package parking_test

import (
	"context"
	"testing"

	"github.com/fjrid/parking/internal/app/model/constant"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/internal/app/service/parking"
	"github.com/fjrid/parking/internal/generated/mock/app/repo/postgresql/parking_mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type OnMockParkingRepo func(mockRepo *parking_mock.MockParkingRepository)

func createParkingSvc(t *testing.T, onMockParkingRepo OnMockParkingRepo) (*parking.ParkingServiceImpl, *gomock.Controller) {
	mock := gomock.NewController(t)
	mockRepo := parking_mock.NewMockParkingRepository(mock)
	if onMockParkingRepo != nil {
		onMockParkingRepo(mockRepo)
	}

	return &parking.ParkingServiceImpl{
		ParkingRepo: mockRepo,
	}, mock
}

func TestParkingSvc_CreateParking(t *testing.T) {
	testcases := []struct {
		testName          string
		OnMockParkingRepo OnMockParkingRepo
		request           dto.CreateParkingRequest
		expectedErr       string
	}{
		{
			testName: "name is missing",
			request: dto.CreateParkingRequest{
				TotalSlot: 10,
			},
			expectedErr: "name: non zero value required",
		},
		{
			testName: "total slot is missing",
			request: dto.CreateParkingRequest{
				Name: "Parking 1",
			},
			expectedErr: "total_slot: non zero value required",
		},
		{
			testName: "parking name already used",
			request: dto.CreateParkingRequest{
				Name:      "Parking 1",
				TotalSlot: 10,
			},
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindParkingByName(gomock.Any(), gomock.Any()).Return([]entity.Parking{
					entity.Parking{
						Name: "Parking 1",
					},
				}, nil)
			},
			expectedErr: "parking name already used",
		},
		{
			testName: "Success",
			request: dto.CreateParkingRequest{
				Name:      "Parking 1",
				TotalSlot: 1,
			},
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindParkingByName(gomock.Any(), gomock.Any()).Return(make([]entity.Parking, 0), nil)
				mockRepo.EXPECT().CreateParking(gomock.Any(), gomock.Any()).Return(entity.Parking{
					ID: 1,
				}, nil)
				mockRepo.EXPECT().CreateParkingSlot(gomock.Any(), gomock.Any()).Return(entity.ParkingSlot{}, nil)
			},
			expectedErr: "",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createParkingSvc(t, tt.OnMockParkingRepo)
			defer mock.Finish()

			err := svc.CreateParking(context.Background(), tt.request)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParkingSvc_BookingPark(t *testing.T) {
	testcases := []struct {
		testName          string
		OnMockParkingRepo OnMockParkingRepo
		request           dto.BookingParkingRequest
		expectedErr       string
	}{
		{
			testName: "car number is missing",
			request: dto.BookingParkingRequest{
				ParkingSlotID: 1,
			},
			expectedErr: "car_number: non zero value required",
		},
		{
			testName: "slot not found",
			request: dto.BookingParkingRequest{
				ParkingSlotID: 1,
				CarNumber:     "N12314Z",
			},
			expectedErr: "parking slot not found",
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindSlotByID(gomock.Any(), gomock.Any()).Return(make([]entity.ParkingSlot, 0), nil)
			},
		},
		{
			testName: "parking slot unavailable",
			request: dto.BookingParkingRequest{
				ParkingSlotID: 1,
				CarNumber:     "N12314Z",
			},
			expectedErr: "parking slot unavailable",
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindSlotByID(gomock.Any(), gomock.Any()).Return([]entity.ParkingSlot{
					entity.ParkingSlot{
						ID:     1,
						Status: constant.PARKING_STATUS_MAINTENANCE,
					},
				}, nil)
			},
		},
		{
			testName: "user has active booking",
			request: dto.BookingParkingRequest{
				ParkingSlotID: 1,
				CarNumber:     "N12314Z",
			},
			expectedErr: "you have active booking",
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindSlotByID(gomock.Any(), gomock.Any()).Return([]entity.ParkingSlot{
					entity.ParkingSlot{
						ID:     1,
						Status: constant.PARKING_STATUS_AVAILABLE,
					},
				}, nil)
				mockRepo.EXPECT().FindBookByUserIDAndStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.ParkingBook{
					entity.ParkingBook{
						ID: 1,
					},
				}, nil)
			},
		},
		{
			testName: "Success",
			request: dto.BookingParkingRequest{
				ParkingSlotID: 1,
				CarNumber:     "N12314Z",
			},
			expectedErr: "",
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindSlotByID(gomock.Any(), gomock.Any()).Return([]entity.ParkingSlot{
					entity.ParkingSlot{
						ID:     1,
						Status: constant.PARKING_STATUS_AVAILABLE,
					},
				}, nil)
				mockRepo.EXPECT().FindBookByUserIDAndStatus(gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]entity.ParkingBook, 0), nil)
				mockRepo.EXPECT().CreateParkingBook(gomock.Any(), gomock.Any()).Return(entity.ParkingBook{}, nil)
				mockRepo.EXPECT().UpdateParkingSlot(gomock.Any(), gomock.Any()).Return(entity.ParkingSlot{}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createParkingSvc(t, tt.OnMockParkingRepo)
			defer mock.Finish()

			_, err := svc.BookingPark(context.Background(), tt.request)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParkingSvc_FinishBook(t *testing.T) {
	testcases := []struct {
		testName          string
		OnMockParkingRepo OnMockParkingRepo
		request           dto.FinishBookRequest
		expectedErr       string
	}{
		{
			testName: "booking not found",
			request: dto.FinishBookRequest{
				BookID: 1,
			},
			expectedErr: "booking not found",
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindBookByIDAndUserIDAndStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]entity.ParkingBook, 0), nil)
			},
		},
		{
			testName: "parking slot not found",
			request: dto.FinishBookRequest{
				BookID: 1,
			},
			expectedErr: "parking slot not found",
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindBookByIDAndUserIDAndStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.ParkingBook{
					entity.ParkingBook{
						ID: 1,
					},
				}, nil)
				mockRepo.EXPECT().UpdateParkingBook(gomock.Any(), gomock.Any()).Return(entity.ParkingBook{}, nil)
				mockRepo.EXPECT().FindSlotByID(gomock.Any(), gomock.Any()).Return(make([]entity.ParkingSlot, 0), nil)
			},
		},
		{
			testName: "Success",
			request: dto.FinishBookRequest{
				BookID: 1,
			},
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindBookByIDAndUserIDAndStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return([]entity.ParkingBook{
					entity.ParkingBook{
						ID: 1,
					},
				}, nil)
				mockRepo.EXPECT().UpdateParkingBook(gomock.Any(), gomock.Any()).Return(entity.ParkingBook{}, nil)
				mockRepo.EXPECT().FindSlotByID(gomock.Any(), gomock.Any()).Return([]entity.ParkingSlot{
					entity.ParkingSlot{
						ID: 1,
					},
				}, nil)
				mockRepo.EXPECT().UpdateParkingSlot(gomock.Any(), gomock.Any()).Return(entity.ParkingSlot{}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createParkingSvc(t, tt.OnMockParkingRepo)
			defer mock.Finish()

			_, err := svc.FinishBook(context.Background(), tt.request)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParkingSvc_ListSlot(t *testing.T) {
	testcases := []struct {
		testName          string
		OnMockParkingRepo OnMockParkingRepo
		viewPagination    dto.ViewPagination
		request           dto.ListSlotParkRequest
		expectedErr       string
	}{
		{
			testName: "Success",
			request: dto.ListSlotParkRequest{
				Status:     "AVAILABLE",
				CarNumber:  "N12314ASD",
				ParkName:   "Parking 1",
				ParkNumber: "1",
			},
			viewPagination: dto.ViewPagination{},
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().FindParkingSlot(gomock.Any(), gomock.Any(), gomock.Any()).Return(make([]entity.SummaryParkingSlot, 0), dto.ViewPagination{}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createParkingSvc(t, tt.OnMockParkingRepo)
			defer mock.Finish()

			_, _, err := svc.ListSlotPark(context.Background(), tt.viewPagination, tt.request)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParkingSvc_ParkingBookSummary(t *testing.T) {
	testcases := []struct {
		testName          string
		OnMockParkingRepo OnMockParkingRepo
		request           dto.ParkingBookSummaryRequest
		expectedErr       string
	}{
		{
			testName: "Start Date Empty",
			request: dto.ParkingBookSummaryRequest{
				EndDate: "2023-01-02",
			},
			expectedErr: "start_date: non zero value required",
		},
		{
			testName: "End Date Empty",
			request: dto.ParkingBookSummaryRequest{
				StartDate: "2023-01-02",
			},
			expectedErr: "end_date: non zero value required",
		},
		{
			testName: "Success",
			request: dto.ParkingBookSummaryRequest{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-02",
			},
			OnMockParkingRepo: func(mockRepo *parking_mock.MockParkingRepository) {
				mockRepo.EXPECT().ParkingBookSummary(gomock.Any(), gomock.Any(), gomock.Any()).Return(dto.SummaryParkingBook{}, nil)
			},
		},
	}

	for _, tt := range testcases {
		t.Run(tt.testName, func(t *testing.T) {
			svc, mock := createParkingSvc(t, tt.OnMockParkingRepo)
			defer mock.Finish()

			_, err := svc.ParkingBookSummary(context.Background(), tt.request)

			if tt.expectedErr != "" {
				require.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
