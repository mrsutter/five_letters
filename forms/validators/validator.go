package validators

import (
  "five_letters/models"
)

type Validator interface {
  Validate() models.Error
}
