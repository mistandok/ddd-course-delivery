package get_all_couriers

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

type GetAllCouriersResponse struct {
	Couriers []CourierDTO
}

type CourierDTO struct {
	ID       uuid.UUID   `db:"id"`
	Name     string      `db:"name"`
	Location LocationDTO `db:"location"`
}

type LocationDTO struct {
	X int64
	Y int64
}

func (l *LocationDTO) String() string {
	return fmt.Sprintf("(%d,%d)", l.X, l.Y)
}

func (l *LocationDTO) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		b, ok := src.([]byte)
		if !ok {
			return errors.New("не удалось преобразовать POINT")
		}
		s = string(b)
	}

	re, err := regexp.Compile(`\((-?\d+\.?\d*),(-?\d+\.?\d*)\)`)
	if err != nil {
		return err
	}

	parts := re.FindStringSubmatch(s)
	if len(parts) != 3 {
		return fmt.Errorf("неожиданный формат POINT: %q", s)
	}
	x, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return err
	}
	y, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return err
	}

	l.X, l.Y = x, y

	return nil
}
