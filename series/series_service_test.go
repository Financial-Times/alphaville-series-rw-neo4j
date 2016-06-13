package series

import (
	"os"
	"testing"

	"github.com/Financial-Times/base-ft-rw-app-go/baseftrwapp"
	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/jmcvetta/neoism"
	"github.com/stretchr/testify/assert"
)

var seriesDriver baseftrwapp.Service

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"

	seriesDriver = getSeriesCypherDriver(t)

	seriesToDelete := Series{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(seriesDriver.Write(seriesToDelete), "Failed to write series")

	found, err := seriesDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete series for uuid %", uuid)
	assert.NoError(err, "Error deleting series for uuid %s", uuid)

	p, found, err := seriesDriver.Read(uuid)

	assert.Equal(Series{}, p, "Found series %s who should have been deleted", p)
	assert.False(found, "Found series for uuid %s who should have been deleted", uuid)
	assert.NoError(err, "Error trying to find series for uuid %s", uuid)
}

func TestCreateAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	seriesDriver = getSeriesCypherDriver(t)

	seriesToWrite := Series{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")

	readSeriesForUUIDAndCheckFieldsMatch(t, uuid, seriesToWrite)

	cleanUp(t, uuid)
}

func TestCreateHandlesSpecialCharacters(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	seriesDriver = getSeriesCypherDriver(t)

	seriesToWrite := Series{UUID: uuid, PrefLabel: "Test 'special chars", TmeIdentifier: "TME_ID"}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")

	readSeriesForUUIDAndCheckFieldsMatch(t, uuid, seriesToWrite)

	cleanUp(t, uuid)
}

func TestCreateNotAllValuesPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	seriesDriver = getSeriesCypherDriver(t)

	seriesToWrite := Series{UUID: uuid, PrefLabel: "Test"}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")

	readSeriesForUUIDAndCheckFieldsMatch(t, uuid, seriesToWrite)

	cleanUp(t, uuid)
}

func TestUpdateWillRemovePropertiesNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	uuid := "12345"
	seriesDriver = getSeriesCypherDriver(t)

	seriesToWrite := Series{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")
	readSeriesForUUIDAndCheckFieldsMatch(t, uuid, seriesToWrite)

	updatedSeries := Series{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(seriesDriver.Write(updatedSeries), "Failed to write updated series")
	readSeriesForUUIDAndCheckFieldsMatch(t, uuid, updatedSeries)

	cleanUp(t, uuid)
}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	seriesDriver = getSeriesCypherDriver(t)
	err := seriesDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
}

func getSeriesCypherDriver(t *testing.T) service {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	db, err := neoism.Connect(url)
	assert.NoError(err, "Failed to connect to Neo4j")
	return NewCypherSeriesService(neoutils.StringerDb{db}, db)
}

func readSeriesForUUIDAndCheckFieldsMatch(t *testing.T, uuid string, expectedSeries Series) {
	assert := assert.New(t)
	storedSeries, found, err := seriesDriver.Read(uuid)

	assert.NoError(err, "Error finding series for uuid %s", uuid)
	assert.True(found, "Didn't find series for uuid %s", uuid)
	assert.Equal(expectedSeries, storedSeries, "series should be the same")
}

func TestWritePrefLabelIsAlsoWrittenAndIsEqualToName(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getSeriesCypherDriver(t)
	uuid := "12345"
	seriesToWrite := Series{UUID: uuid, PrefLabel: "Test", TmeIdentifier: "TME_ID"}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")

	result := []struct {
		PrefLabel string `json:"t.prefLabel"`
	}{}

	getPrefLabelQuery := &neoism.CypherQuery{
		Statement: `
				MATCH (t:Series {uuid:"12345"}) RETURN t.prefLabel
				`,
		Result: &result,
	}

	err := seriesDriver.cypherRunner.CypherBatch([]*neoism.CypherQuery{getPrefLabelQuery})
	assert.NoError(err)
	assert.Equal("Test", result[0].PrefLabel, "PrefLabel should be 'Test")
	cleanUp(t, uuid)
}

func cleanUp(t *testing.T, uuid string) {
	assert := assert.New(t)
	found, err := seriesDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete series for uuid %", uuid)
	assert.NoError(err, "Error deleting series for uuid %s", uuid)
}
