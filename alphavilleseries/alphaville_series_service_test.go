package alphavilleseries

import (
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

var alphavilleSeriesDriver baseftrwapp.Service

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"

	alphavilleSeriesDriver = getAlphavilleSeriesCypherDriver(t)

	alphavilleSeriesToDelete := AlphavilleSeries{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(alphavilleSeriesDriver.Write(alphavilleSeriesToDelete), "Failed to write Alphaville Series")

	found, err := alphavilleSeriesDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete Alphaville Series for uuid %s", uuid)
	assert.NoError(err, "Error deleting Alphaville Series for uuid %s", uuid)

	p, found, err := alphavilleSeriesDriver.Read(uuid)

	assert.Equal(AlphavilleSeries{}, p, "Found Alphaville Series %s which should have been deleted", p)
	assert.False(found, "Found Alphaville Series for uuid %s which should have been deleted", uuid)
	assert.NoError(err, "Error trying to find Alphaville Series for uuid %s", uuid)
}

func TestCreateAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	alphavilleSeriesDriver = getAlphavilleSeriesCypherDriver(t)

	alphavilleSeriesToWrite := AlphavilleSeries{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(alphavilleSeriesDriver.Write(alphavilleSeriesToWrite), "Failed to write Alphaville Series")

	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(t, uuid, alphavilleSeriesToWrite)

	cleanUp(t, uuid)
}

func TestCreateHandlesSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	alphavilleSeriesDriver = getAlphavilleSeriesCypherDriver(t)

	alphavilleSeriesToWrite := AlphavilleSeries{UUID: uuid, PrefLabel: "Test 'special chars", TmeIdentifier: "TME_ID"}

	assert.NoError(alphavilleSeriesDriver.Write(alphavilleSeriesToWrite), "Failed to write Alphaville Series")

	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(t, uuid, alphavilleSeriesToWrite)

	cleanUp(t, uuid)
}

func TestCreateNotAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	alphavilleSeriesDriver = getAlphavilleSeriesCypherDriver(t)

	alphavilleSeriesToWrite := AlphavilleSeries{UUID: uuid, PrefLabel: "Test"}

	assert.NoError(alphavilleSeriesDriver.Write(alphavilleSeriesToWrite), "Failed to write Alphaville Series")

	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(t, uuid, alphavilleSeriesToWrite)

	cleanUp(t, uuid)
}

func TestUpdateWillRemovePropertiesNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	alphavilleSeriesDriver = getAlphavilleSeriesCypherDriver(t)

	alphavilleSeriesToWrite := AlphavilleSeries{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(alphavilleSeriesDriver.Write(alphavilleSeriesToWrite), "Failed to write Alphaville Series")
	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(t, uuid, alphavilleSeriesToWrite)

	updatedAlphavilleSeries := AlphavilleSeries{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(alphavilleSeriesDriver.Write(updatedAlphavilleSeries), "Failed to write updated Alphaville Series")
	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(t, uuid, updatedAlphavilleSeries)

	cleanUp(t, uuid)
}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	alphavilleSeriesDriver = getAlphavilleSeriesCypherDriver(t)
	err := alphavilleSeriesDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
}

func getAlphavilleSeriesCypherDriver(t *testing.T) service {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(url)
	assert.NoError(err, "Failed to connect to Neo4j")
	return NewCypherAlphavilleSeriesService(neoutils.StringerDb{db}, db)
}

func readAlphavilleSeriesForUUIDAndCheckFieldsMatch(t *testing.T, uuid string, expectedAlphavilleSeries AlphavilleSeries) {
	assert := assert.New(t)
	storedAlphavilleSeries, found, err := alphavilleSeriesDriver.Read(uuid)

	assert.NoError(err, "Error finding Alphaville Series for uuid %s", uuid)
	assert.True(found, "Didn't find Alphaville Series for uuid %s", uuid)
	assert.Equal(expectedAlphavilleSeries, storedAlphavilleSeries, "Alphaville Series should be the same")
}

func TestWritePrefLabelIsAlsoWrittenAndIsEqualToName(t *testing.T) {
	assert := assert.New(t)
	alphavilleSeriesDriver := getAlphavilleSeriesCypherDriver(t)
	uuid := "12345"
	alphavilleSeriesToWrite := AlphavilleSeries{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(alphavilleSeriesDriver.Write(alphavilleSeriesToWrite), "Failed to write Alphaville Series")

	result := []struct {
		PrefLabel string `json:"t.prefLabel"`
	}{}

	getPrefLabelQuery := &neoism.CypherQuery{
		Statement: `
				MATCH (t:AlphavilleSeries {uuid:"12345"}) RETURN t.prefLabel
				`,
		Result: &result,
	}

	err := alphavilleSeriesDriver.cypherRunner.CypherBatch([]*neoism.CypherQuery{getPrefLabelQuery})
	assert.NoError(err)
	assert.Equal("Test", result[0].PrefLabel, "PrefLabel should be 'Test")
	cleanUp(t, uuid)
}

func cleanUp(t *testing.T, uuid string) {
	assert := assert.New(t)
	found, err := alphavilleSeriesDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete Alphaville Series for uuid %s", uuid)
	assert.NoError(err, "Error deleting Alphaville Series for uuid %s", uuid)
}
