package parse

import (
	"fmt"
	"strings"
)

type SynapseRoleAssignmentId struct {
	WorkspaceName string
	Id            string
}

func SynapseRoleAssignmentID(input string) (*SynapseRoleAssignmentId, error) {
	segments := strings.Split(input, "|")
	if len(segments) != 2 {
		return nil, fmt.Errorf("expected an ID in the format `{workspaceName}|{id} but got %q", input)
	}

	return &SynapseRoleAssignmentId{
		WorkspaceName: segments[0],
		Id:            segments[1],
	}, nil
}
