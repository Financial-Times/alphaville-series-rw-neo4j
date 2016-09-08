package alphavilleseries

import (
	"os"
	"testing"

	"github.com/Financial-Times/neo-utils-go/neoutils"
	"github.com/stretchr/testify/assert"
)

const (
	seriesUUID              = "12345"
	newAlphavilleSeriesUUID = "123456"
	tmeID                   = "TME_ID"
	newTmeID                = "NEW_TME_ID"
	prefLabel               = "Test"
	specialCharPrefLabel    = "Test 'special chars"
)

var defaultTypes = []string{"Thing", "Concept", "Classification", "AlphavilleSeries"}

func TestConnectivityCheck(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)
	err := seriesDriver.Check()
	assert.NoError(err, "Unexpected error on connectivity check")
}

func TestPrefLabelIsCorrectlyWritten(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)

	alternativeIdentifiers := alternativeIdentifiers{UUIDS: []string{seriesUUID}}
	seriesToWrite := AlphavilleSeries{UUID: seriesUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	err := seriesDriver.Write(seriesToWrite)
	assert.NoError(err, "ERROR happened during write time")

	storedAlphavilleSeries, found, err := seriesDriver.Read(seriesUUID)
	assert.NoError(err, "ERROR happened during read time")
	assert.Equal(true, found)
	assert.NotEmpty(storedAlphavilleSeries)

	assert.Equal(prefLabel, storedAlphavilleSeries.(AlphavilleSeries).PrefLabel, "PrefLabel should be "+prefLabel)
	cleanUp(assert, seriesUUID, seriesDriver)
}

func TestPrefLabelSpecialCharactersAreHandledByCreate(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{seriesUUID}}
	seriesToWrite := AlphavilleSeries{UUID: seriesUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")

	//add default types that will be automatically added by the writer
	seriesToWrite.Types = defaultTypes
	//check if seriesToWrite is the same with the one inside the DB
	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(assert, seriesDriver, seriesUUID, seriesToWrite)
	cleanUp(assert, seriesUUID, seriesDriver)
}

func TestCreateCompleteAlphavilleSeriesWithPropsAndIdentifiers(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{seriesUUID}}
	seriesToWrite := AlphavilleSeries{UUID: seriesUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")

	//add default types that will be automatically added by the writer
	seriesToWrite.Types = defaultTypes
	//check if seriesToWrite is the same with the one inside the DB
	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(assert, seriesDriver, seriesUUID, seriesToWrite)
	cleanUp(assert, seriesUUID, seriesDriver)
}

func TestUpdateWillRemovePropertiesAndIdentifiersNoLongerPresent(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)

	allAlternativeIdentifiers := alternativeIdentifiers{TME: []string{}, UUIDS: []string{seriesUUID}}
	seriesToWrite := AlphavilleSeries{UUID: seriesUUID, PrefLabel: prefLabel, AlternativeIdentifiers: allAlternativeIdentifiers}

	assert.NoError(seriesDriver.Write(seriesToWrite), "Failed to write series")
	//add default types that will be automatically added by the writer
	seriesToWrite.Types = defaultTypes
	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(assert, seriesDriver, seriesUUID, seriesToWrite)

	tmeAlternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{seriesUUID}}
	updatedAlphavilleSeries := AlphavilleSeries{UUID: seriesUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: tmeAlternativeIdentifiers}

	assert.NoError(seriesDriver.Write(updatedAlphavilleSeries), "Failed to write updated series")
	//add default types that will be automatically added by the writer
	updatedAlphavilleSeries.Types = defaultTypes
	readAlphavilleSeriesForUUIDAndCheckFieldsMatch(assert, seriesDriver, seriesUUID, updatedAlphavilleSeries)

	cleanUp(assert, seriesUUID, seriesDriver)
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)

	alternativeIdentifiers := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{seriesUUID}}
	seriesToDelete := AlphavilleSeries{UUID: seriesUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIdentifiers}

	assert.NoError(seriesDriver.Write(seriesToDelete), "Failed to write series")

	found, err := seriesDriver.Delete(seriesUUID)
	assert.True(found, "Didn't manage to delete series for uuid %s", seriesUUID)
	assert.NoError(err, "Error deleting series for uuid %s", seriesUUID)

	p, found, err := seriesDriver.Read(seriesUUID)

	assert.Equal(AlphavilleSeries{}, p, "Found series %s who should have been deleted", p)
	assert.False(found, "Found series for uuid %s who should have been deleted", seriesUUID)
	assert.NoError(err, "Error trying to find series for uuid %s", seriesUUID)
}

func TestCount(t *testing.T) {
	assert := assert.New(t)
	seriesDriver := getAlphavilleSeriesService(t)

	alternativeIds := alternativeIdentifiers{TME: []string{tmeID}, UUIDS: []string{seriesUUID}}
	seriesOneToCount := AlphavilleSeries{UUID: seriesUUID, PrefLabel: prefLabel, AlternativeIdentifiers: alternativeIds}

	assert.NoError(seriesDriver.Write(seriesOneToCount), "Failed to write series")

	nr, err := seriesDriver.Count()
	assert.Equal(1, nr, "Should be 1 series in DB - count differs")
	assert.NoError(err, "An unexpected error occurred during count")

	newAlternativeIds := alternativeIdentifiers{TME: []string{newTmeID}, UUIDS: []string{newAlphavilleSeriesUUID}}
	seriesTwoToCount := AlphavilleSeries{UUID: newAlphavilleSeriesUUID, PrefLabel: specialCharPrefLabel, AlternativeIdentifiers: newAlternativeIds}

	assert.NoError(seriesDriver.Write(seriesTwoToCount), "Failed to write series")

	nr, err = seriesDriver.Count()
	assert.Equal(2, nr, "Should be 2 series in DB - count differs")
	assert.NoError(err, "An unexpected error occurred during count")

	cleanUp(assert, seriesUUID, seriesDriver)
	cleanUp(assert, newAlphavilleSeriesUUID, seriesDriver)
}

func readAlphavilleSeriesForUUIDAndCheckFieldsMatch(assert *assert.Assertions, seriesDriver service, uuid string, expectedAlphavilleSeries AlphavilleSeries) {

	storedAlphavilleSeries, found, err := seriesDriver.Read(uuid)
	assert.NoError(err, "Error finding series for uuid %s", uuid)
	assert.True(found, "Didn't find series for uuid %s", uuid)
	assert.Equal(expectedAlphavilleSeries.PrefLabel, storedAlphavilleSeries.(AlphavilleSeries).PrefLabel, "PerfLable didn't match")
	assert.Equal(expectedAlphavilleSeries.AlternativeIdentifiers, storedAlphavilleSeries.(AlphavilleSeries).AlternativeIdentifiers, "AlternativeIdentifiers didn't match")
	assert.True(sameSetOfValues(expectedAlphavilleSeries.Types, storedAlphavilleSeries.(AlphavilleSeries).Types), "Types should be the same")
}

func getAlphavilleSeriesService(t *testing.T) service {
	assert := assert.New(t)
	url := os.Getenv("NEO4J_TEST_URL")
	if url == "" {
		url = "http://localhost:7474/db/data"
	}

	conf := neoutils.DefaultConnectionConfig()
	conf.Transactional = false
	db, err := neoutils.Connect(url, conf)
	assert.NoError(err, "Failed to connect to Neo4j")

	service := NewCypherAlphavilleSeriesService(db)
	service.Initialise()
	return service
}

func cleanUp(assert *assert.Assertions, uuid string, seriesDriver service) {
	found, err := seriesDriver.Delete(uuid)
	assert.True(found, "Didn't manage to delete series for uuid %s", uuid)
	assert.NoError(err, "Error deleting series for uuid %s", uuid)
}

func TestSameSetOfValues(t *testing.T) {
	x := []string{"1", "2"}
	y := []string{"2", "1"}

	assert := assert.New(t)
	assert.True(sameSetOfValues(x, y), "Arrays were equial but evaluated as not equal")
}

func sameSetOfValues(x []string, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	m := make(map[string]int)

	for _, x_value := range x {
		m[x_value]++
	}
	for _, y_value := range y {
		if m[y_value] == 0 {
			return false
		}
	}
	return true
}
