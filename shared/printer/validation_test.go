package printer

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestValidatePrintRequest(t *testing.T) {
    validItems := []PrintItem{
        {Type: Text, Content: "Hello", Align: AlignLeft},
        {Type: QRCode, Content: "QR data", Align: AlignCenter},
        {Type: Blank, Count: 1},
        {Type: Line},
        {Type: Cut},
    }

    validRequest := PrintRequest{
        Items: validItems,
    }

    // Valid request should pass
    err := ValidatePrintRequest(validRequest)
    assert.NoError(t, err)

    // Test invalid item type
    invalidTypeItems := []PrintItem{
        {Type: ItemType("invalid")},
    }
    invalidRequest := PrintRequest{Items: invalidTypeItems}
    err = ValidatePrintRequest(invalidRequest)
    assert.Error(t, err)

    // Test invalid alignment
    invalidAlignItems := []PrintItem{
        {Type: Text, Align: Alignment("invalid")},
    }
    invalidAlignRequest := PrintRequest{Items: invalidAlignItems}
    err = ValidatePrintRequest(invalidAlignRequest)
    assert.Error(t, err)

    // Test missing required type
    missingTypeItems := []PrintItem{
        {Content: "No type"},
    }
    missingTypeRequest := PrintRequest{Items: missingTypeItems}
    err = ValidatePrintRequest(missingTypeRequest)
    assert.Error(t, err)

    // Test empty items list (should fail because of required)
    emptyRequest := PrintRequest{Items: []PrintItem{}}
    err = ValidatePrintRequest(emptyRequest)
    assert.Error(t, err)
}
