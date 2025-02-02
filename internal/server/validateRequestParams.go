package server

import (
	"errors"
	"strconv"

	"github.com/EvansTrein/iqProgers/models"
)

func validateRequestParams(params map[string]string, reqStruct interface{}) error {
	switch req := reqStruct.(type) {
	case *models.UserOperationsRequest:
		return validateUserOperationsRequest(params, req)
	default:
		return errors.New("unsupported request type")
	}
}

func validateUserOperationsRequest(params map[string]string, req *models.UserOperationsRequest) error {
	userId, ok := params["userID"]
	if !ok {
		return errors.New("no user id in params")
	}

	offset, ok := params["offset"]
	if !ok {
		return errors.New("no offset in params")
	}

	limit, ok := params["limit"]
	if !ok {
		return errors.New("no limit in params")
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}

	offsetInt, err := strconv.Atoi(offset)
	if err != nil {
		return err
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		return err
	}

	req.UserID = uint(userIdInt)
	req.Offset = offsetInt
	req.Limit = limitInt

	return nil
}
