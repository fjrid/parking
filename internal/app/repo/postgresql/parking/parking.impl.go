package parking

import (
	"context"
	"fmt"
	"time"

	"github.com/fjrid/parking/internal/app/model/constant"
	"github.com/fjrid/parking/internal/app/model/dto"
	"github.com/fjrid/parking/internal/app/model/entity"
	"github.com/fjrid/parking/pkg/dbtxn"
	"github.com/fjrid/parking/pkg/sqkit"

	sq "github.com/Masterminds/squirrel"
)

// @ctor
func NewParkingRepository(impl ParkingRepositoryImpl) ParkingRepository {
	return &impl
}

func (r *ParkingRepositoryImpl) CreateParking(ctx context.Context, data dto.CreateParkingRequest) (res entity.Parking, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	res = entity.Parking{
		Name:       data.Name,
		CreatedAt:  time.Now(),
		CreatedBy:  data.CreatedBy,
		ModifiedAt: time.Now(),
		ModifiedBy: data.CreatedBy,
	}

	builder := sq.
		Insert(entity.ParkingTableName).
		Columns(
			entity.ParkingTable.Name,
			entity.ParkingTable.CreatedBy,
			entity.ParkingTable.ModifiedBy,
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", entity.ParkingTable.ID),
		).
		PlaceholderFormat(sq.Dollar).
		Values(
			res.Name,
			res.CreatedBy,
			res.ModifiedBy,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	if err = scanner.Scan(&res.ID); err != nil {
		txn.AppendError(err)
		return
	}

	return
}

func (r *ParkingRepositoryImpl) CreateParkingSlot(ctx context.Context, data dto.CreateParkingSlotRequest) (res entity.ParkingSlot, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	res = entity.ParkingSlot{
		ParkingID:  data.ParkingID,
		Number:     data.Number,
		Status:     data.Status,
		CreatedAt:  time.Now(),
		CreatedBy:  data.CreatedBy,
		ModifiedAt: time.Now(),
		ModifiedBy: data.CreatedBy,
	}

	builder := sq.
		Insert(entity.ParkingSlotTableName).
		Columns(
			entity.ParkingSlotTable.ParkingID,
			entity.ParkingSlotTable.Number,
			entity.ParkingSlotTable.Status,
			entity.ParkingSlotTable.CreatedBy,
			entity.ParkingSlotTable.ModifiedBy,
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", entity.ParkingSlotTable.ID),
		).
		PlaceholderFormat(sq.Dollar).
		Values(
			res.ParkingID,
			res.Number,
			res.Status,
			res.CreatedBy,
			res.ModifiedBy,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	if err = scanner.Scan(&res.ID); err != nil {
		txn.AppendError(err)
		return
	}

	return
}

func (r *ParkingRepositoryImpl) UpdateParkingSlot(ctx context.Context, data entity.ParkingSlot) (res entity.ParkingSlot, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.
		Update(entity.ParkingSlotTableName).
		Set(entity.ParkingSlotTable.Status, data.Status).
		Set(entity.ParkingSlotTable.ModifiedAt, data.ModifiedAt).
		Set(entity.ParkingSlotTable.ModifiedBy, data.ModifiedBy).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn).
		Where(fmt.Sprintf("%s.%s = %d", entity.ParkingSlotTableName, entity.ParkingSlotTable.ID, data.ID))

	_, err = builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return
	}

	return
}

func (r *ParkingRepositoryImpl) FindParking(ctx context.Context, opts ...sqkit.SelectOption) (results []entity.Parking, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.
		Select(
			"*",
		).
		From(entity.ParkingTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	results = make([]entity.Parking, 0)
	for rows.Next() {
		parking := entity.Parking{}

		if err = rows.Scan(
			&parking.ID,
			&parking.Name,
			&parking.CreatedAt,
			&parking.CreatedBy,
			&parking.ModifiedAt,
			&parking.ModifiedBy,
			&parking.DeletedAt,
			&parking.DeletedBy,
		); err != nil {
			return
		}

		results = append(results, parking)
	}

	return
}

func (r *ParkingRepositoryImpl) FindParkingByName(ctx context.Context, name string) (results []entity.Parking, err error) {
	condition := []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingTableName, entity.ParkingTable.Name): name},
	}

	return r.FindParking(ctx, condition...)
}

func (r *ParkingRepositoryImpl) FindSlot(ctx context.Context, opts ...sqkit.SelectOption) (results []entity.ParkingSlot, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.
		Select(
			fmt.Sprintf("%s.*", entity.ParkingSlotTableName),
			fmt.Sprintf("%s.%s", entity.ParkingTableName, entity.ParkingTable.Name),
		).
		From(entity.ParkingSlotTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn).
		Join(fmt.Sprintf("%s ON %s.%s = %s.%s", entity.ParkingTableName, entity.ParkingSlotTableName, entity.ParkingSlotTable.ParkingID, entity.ParkingTableName, entity.ParkingTable.ID))

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	results = make([]entity.ParkingSlot, 0)
	for rows.Next() {
		parkingSlot := entity.ParkingSlot{}

		if err = rows.Scan(
			&parkingSlot.ID,
			&parkingSlot.ParkingID,
			&parkingSlot.Number,
			&parkingSlot.Status,
			&parkingSlot.CreatedAt,
			&parkingSlot.CreatedBy,
			&parkingSlot.ModifiedAt,
			&parkingSlot.ModifiedBy,
			&parkingSlot.DeletedAt,
			&parkingSlot.DeletedBy,
			&parkingSlot.ParkingName,
		); err != nil {
			return
		}

		results = append(results, parkingSlot)
	}

	return
}

func (r *ParkingRepositoryImpl) FindSlotByID(ctx context.Context, id int64) (results []entity.ParkingSlot, err error) {
	condition := []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.ID): id},
	}

	return r.FindSlot(ctx, condition...)
}

func (r *ParkingRepositoryImpl) CreateParkingBook(ctx context.Context, data dto.CreateParkingBookRequest) (res entity.ParkingBook, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	res = entity.ParkingBook{
		UserID:        data.UserID,
		ParkingSlotID: data.ParkingSlotID,
		StartTime:     data.StartTime,
		Status:        data.Status,
		CarNumber:     data.CarNumber,
		CreatedAt:     time.Now(),
		CreatedBy:     data.CreatedBy,
		ModifiedAt:    time.Now(),
		ModifiedBy:    data.CreatedBy,
	}

	builder := sq.
		Insert(entity.ParkingBookTableName).
		Columns(
			entity.ParkingBookTable.UserID,
			entity.ParkingBookTable.ParkingSlotID,
			entity.ParkingBookTable.StartTime,
			entity.ParkingBookTable.Status,
			entity.ParkingBookTable.CarNumber,
			entity.ParkingBookTable.CreatedBy,
			entity.ParkingBookTable.ModifiedBy,
		).
		Suffix(
			fmt.Sprintf("RETURNING \"%s\"", entity.ParkingTable.ID),
		).
		PlaceholderFormat(sq.Dollar).
		Values(
			res.UserID,
			res.ParkingSlotID,
			res.StartTime,
			res.Status,
			res.CarNumber,
			res.CreatedBy,
			res.ModifiedBy,
		)

	scanner := builder.RunWith(txn).QueryRowContext(ctx)

	if err = scanner.Scan(&res.ID); err != nil {
		txn.AppendError(err)
		return
	}

	return
}

func (r *ParkingRepositoryImpl) FindBook(ctx context.Context, opts ...sqkit.SelectOption) (results []entity.ParkingBook, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.
		Select(
			"*",
		).
		From(entity.ParkingBookTableName).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn)

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	rows, err := builder.QueryContext(ctx)
	if err != nil {
		return
	}

	results = make([]entity.ParkingBook, 0)
	for rows.Next() {
		parkingBook := entity.ParkingBook{}

		if err = rows.Scan(
			&parkingBook.ID,
			&parkingBook.UserID,
			&parkingBook.ParkingSlotID,
			&parkingBook.StartTime,
			&parkingBook.EndTime,
			&parkingBook.Fee,
			&parkingBook.Status,
			&parkingBook.CarNumber,
			&parkingBook.CreatedAt,
			&parkingBook.CreatedBy,
			&parkingBook.ModifiedAt,
			&parkingBook.ModifiedBy,
			&parkingBook.DeletedAt,
			&parkingBook.DeletedBy,
		); err != nil {
			return
		}

		results = append(results, parkingBook)
	}

	return
}

func (r *ParkingRepositoryImpl) FindBookByUserIDAndStatus(ctx context.Context, userID int64, status string) (results []entity.ParkingBook, err error) {
	condition := []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingBookTableName, entity.ParkingBookTable.UserID): userID},
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingBookTableName, entity.ParkingBookTable.Status): status},
	}

	return r.FindBook(ctx, condition...)
}

func (r *ParkingRepositoryImpl) FindBookByIDAndUserIDAndStatus(ctx context.Context, id int64, userID int64, status string) (results []entity.ParkingBook, err error) {
	condition := []sqkit.SelectOption{
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingBookTableName, entity.ParkingBookTable.ID): id},
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingBookTableName, entity.ParkingBookTable.UserID): userID},
		sqkit.Eq{fmt.Sprintf("%s.%s", entity.ParkingBookTableName, entity.ParkingBookTable.Status): status},
	}

	return r.FindBook(ctx, condition...)
}

func (r *ParkingRepositoryImpl) UpdateParkingBook(ctx context.Context, data entity.ParkingBook) (res entity.ParkingBook, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.
		Update(entity.ParkingBookTableName).
		Set(entity.ParkingBookTable.EndTime, data.EndTime).
		Set(entity.ParkingBookTable.Fee, data.Fee).
		Set(entity.ParkingBookTable.Status, data.Status).
		Set(entity.ParkingBookTable.ModifiedAt, data.ModifiedAt).
		Set(entity.ParkingBookTable.ModifiedBy, data.ModifiedBy).
		PlaceholderFormat(sq.Dollar).
		RunWith(txn).
		Where(fmt.Sprintf("%s.%s = %d", entity.ParkingBookTableName, entity.ParkingBookTable.ID, data.ID))

	_, err = builder.ExecContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return
	}

	return
}

func (r *ParkingRepositoryImpl) FindParkingSlot(ctx context.Context, viewPagination dto.ViewPagination, opts ...sqkit.SelectOption) (results []entity.SummaryParkingSlot, vp dto.ViewPagination, err error) {
	txn, err := dbtxn.Use(ctx, r.DB)
	if err != nil {
		return
	}

	builder := sq.Select(
		"COUNT(parking_slot.id) OVER() AS total",
		fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.ID),
		fmt.Sprintf("%s.%s", entity.ParkingTableName, entity.ParkingTable.Name),
		fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.Number),
		fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.Status),
		fmt.Sprintf("%s.%s", entity.ParkingBookTableName, entity.ParkingBookTable.CarNumber),
	).
		From(entity.ParkingSlotTableName).
		LeftJoin(fmt.Sprintf(
			"%s ON %s.%s = %s.%s AND %s.%s = '%s'",
			entity.ParkingBookTableName, entity.ParkingSlotTableName, entity.ParkingSlotTable.ID,
			entity.ParkingBookTableName, entity.ParkingBookTable.ParkingSlotID,
			entity.ParkingBookTableName, entity.ParkingBookTable.Status, constant.BOOK_STATUS_ON_GOING,
		)).
		LeftJoin(fmt.Sprintf(
			"%s ON %s.%s = %s.%s",
			entity.ParkingTableName, entity.ParkingTableName, entity.ParkingTable.ID,
			entity.ParkingSlotTableName, entity.ParkingSlotTable.ParkingID,
		)).
		PlaceholderFormat(sq.Dollar).
		OrderBy(fmt.Sprintf("%s.%s", entity.ParkingSlotTableName, entity.ParkingSlotTable.ID))

	for _, opt := range opts {
		builder = opt.CompileSelect(builder)
	}

	if viewPagination.Limit > 0 {
		builder = builder.Limit(uint64(viewPagination.Limit))
	}

	if viewPagination.Offset > 0 {
		builder = builder.Offset(uint64(viewPagination.Offset))
	}

	rows, err := builder.RunWith(txn).QueryContext(ctx)
	if err != nil {
		txn.AppendError(err)
		return
	}

	defer rows.Close()

	results = make([]entity.SummaryParkingSlot, 0)

	for rows.Next() {
		var r entity.SummaryParkingSlot

		err = rows.Scan(
			&viewPagination.Total,
			&r.ID,
			&r.ParkingName,
			&r.ParkingNumber,
			&r.Status,
			&r.CarNumber,
		)
		if err != nil {
			txn.AppendError(err)
			return
		}

		results = append(results, r)
	}

	vp = viewPagination

	return
}

func (r *ParkingRepositoryImpl) ParkingBookSummary(ctx context.Context, startDate string, endDate string) (result dto.SummaryParkingBook, err error) {
	query := `
	SELECT
		SUM(d.totalEpoch) AS totalTime,
		COUNT(d.id) AS totalBookParking,
		SUM(d.fee) AS totalFee
	FROM (
		SELECT
			end_time - start_time as totalEpoch,
			*
		FROM parking_book
		WHERE 
			status = 'FINISHED' AND
			start_time::DATE >= '%s' AND
			end_time::DATE <= '%s'
	) d
	`
	r.DB.QueryRowContext(ctx, fmt.Sprintf(query, startDate, endDate)).Scan(
		&result.TotalTime,
		&result.TotalBooking,
		&result.TotalFee,
	)

	return
}
