package transforms

import (
	"github.com/jpleatherland/spacetraders/internal/spec"
	"fmt"
)


func StructureContract(contract spec.ContractResponse) string {
	return fmt.Sprintf("<p>Deadline: %v</p>", contract.Data.Terms.Deadline)
}
