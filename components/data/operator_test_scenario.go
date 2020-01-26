package data

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/workshop/common"
	"github.com/pavlo67/workshop/common/crud"
	"github.com/pavlo67/workshop/common/logger"

	"github.com/pavlo67/workshop/components/tagger"
)

type OperatorTestCase struct {
	Operator
	crud.Cleaner

	ToSave   Item
	ToUpdate Item
}

func TestCases(dataOp Operator, cleanerOp crud.Cleaner) []OperatorTestCase {
	return []OperatorTestCase{
		{
			Operator: dataOp,
			Cleaner:  cleanerOp,
			ToSave: Item{
				ID:      "",
				URL:     "111111",
				Title:   "345456",
				Summary: "6578gj",
				Embedded: []Item{{
					Title:   "56567",
					Summary: "3333333",
					Tags:    []tagger.Tag{{Label: "1"}, {Label: "332343"}},
				}},
				Data: crud.Data{
					TypeKey: "test",
					Content: `{"AAA": "aaa", "BBB": 222}`,
				},
				Tags:      []tagger.Tag{{Label: "1"}, {Label: "333"}},
				OwnerKey:  actorKey,
				ViewerKey: actorKey,
				History: []crud.Action{{
					Key:    crud.CreatedAction,
					DoneAt: time.Time{},
				}},
			},

			ToUpdate: Item{
				URL:     "22222222",
				Title:   "345456rt",
				Summary: "6578eegj",
				Data: crud.Data{
					TypeKey: "test",
					Content: `{"AAA": "awraa", "BBB": 22552}`,
				},
				Tags:      []tagger.Tag{{Label: "1"}, {Label: "333"}},
				OwnerKey:  actorKey,
				ViewerKey: actorKey,
			},
		},
	}
}

// TODO: тест чистки бази
// TODO: test created_at, updated_at
// TODO: test GetOptions

const numRepeats = 3
const toReadI = 0   // must be < numRepeats
const toUpdateI = 1 // must be < numRepeats
const toDeleteI = 2 // must be < numRepeats

const actorKey = "test"

func Compare(t *testing.T, dataOp Operator, readed *Item, expectedItem Item, l logger.Operator) {
	require.NotNil(t, readed)

	l.Infof("to be saved: %#v", expectedItem)
	l.Infof("readed: %#v", readed)

	for i, action := range expectedItem.History {
		expectedItem.History[i].DoneAt = action.DoneAt.UTC()
	}

	expectedDetails := expectedItem.Data.Content
	expectedItem.Data.Content = ""

	readedDetails := readed.Data.Content
	readed.Data.Content = ""

	// TODO!!! check it carefully
	readed.History = nil
	expectedItem.History = nil

	require.Equal(t, &expectedItem, readed)
	require.Equal(t, expectedDetails, readedDetails)

}

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		l.Debug(i)

		var id [numRepeats]common.ID
		var toSave [numRepeats]Item
		// var data Tag

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner.Clean(nil, nil)
		require.NoError(t, err, "what is the error on .Cleaner()?")

		// test Describe ------------------------------------------------------------------------------------

		//description := tc.Description()
		//
		//keyFields := description.PrimaryKeys()
		//
		//if len(keyFields) > 1 {
		//	require.FailNow(t, "too many key fields", keyFields)
		//} else if len(keyFields) < 1 {
		//	keyFields = append(keyFields, "id")
		//}
		//
		////for _, fieldKey := range tc.DescribedFields {
		////	require.NotEmpty(t, description.FieldsArr[fieldKey], "on .Describe(): "+fieldKey+"???")
		////}

		// test Create --------------------------------------------------------------------------------------

		//var uniques, autoUniques []string
		//
		//for _, field := range description.FieldsArr {
		//	key := field.Key
		//	if field.Unique {
		//		if field.AutoUnique {
		//			autoUniques = append(autoUniques, key)
		//		} else {
		//			uniques = append(uniques, key)
		//		}
		//	}
		//}
		//
		//nativeToCreate, err := tc.ItemToNative(tc.ToSave)
		//require.NoError(t, err)

		//if !tc.ExpectedSaveOk {
		//	_, err = tc.Save([]Tag{tc.ToSave}, nil)
		//	require.ErrStr(t, err, "where is an error on .Save()?")
		//	continue
		//}

		for i := 0; i < numRepeats; i++ {
			toSave[i] = tc.ToSave
			//toSave[i].Details = &tc.DetailsToSave
			idI, err := tc.Save(toSave[i], &crud.SaveOptions{ActorKey: actorKey})
			require.NoError(t, err)
			require.NotEmpty(t, idI)
			id[i] = idI
		}

		// test .Read ----------------------------------------------------------------------------------------

		// if !tc.ExpectedReadOk {
		// 	 _, err = tc.Read(id[toReadI], nil)
		//	 require.ErrStr(t, err)
		//	 continue
		// }

		readedSaved, err := tc.Read(id[toReadI], &crud.GetOptions{ActorKey: actorKey})
		require.NoError(t, err)

		toSave[i].ID = id[toReadI]

		Compare(t, tc, readedSaved, toSave[i], l)

		// test .Update & .Read -----------------------------------------------------------------------------------

		// if !tc.ExpectedUpdateOk {
		//	 err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//	 require.ErrStr(t, err, "where is an error on .Update()?")
		//	 continue
		// }

		tc.ToUpdate.ID = id[toUpdateI]
		// tc.ToUpdate.Details = &tc.DetailsToUpdate

		// !!! .History error

		_, err = tc.Save(tc.ToUpdate, &crud.SaveOptions{ActorKey: actorKey})
		require.Error(t, err)

		// !!! corrected .History

		readedToUpdate, err := tc.Read(id[toUpdateI], &crud.GetOptions{ActorKey: actorKey})
		require.NoError(t, err)
		require.NotNil(t, readedToUpdate)

		tc.ToUpdate.History = readedToUpdate.History

		_, err = tc.Save(tc.ToUpdate, &crud.SaveOptions{ActorKey: actorKey})
		require.NoError(t, err)

		readedUpdated, err := tc.Read(id[toUpdateI], &crud.GetOptions{ActorKey: actorKey})
		require.NoError(t, err)

		// tc.ToUpdate.History.CreatedAt = tc.ToSave.History.CreatedAt // unchanged!!!

		Compare(t, tc, readedUpdated, tc.ToUpdate, l)

		// TODO!!!
		//	if !tc.ExcludeUpdateTest {
		//		var uniquesUpdatable []string
		//		for _, field := range description.FieldsArr {
		//			if field.Unique && (field.Updatable && !field.AutoUnique) { // || field.Additable
		//				uniquesUpdatable = append(uniquesUpdatable, field.Key)
		//			}
		//		}
		//
		//		//tc.ToUpdate[keyFields[0]] = id[toUpdateI]
		//
		//		nativeToUpdate, err := tc.ItemToNative(tc.ToUpdate)
		//		require.NoError(t, err)
		//
		//
		//		if tc.ISToUpdateBad != nil {
		//			err = tc.Update(*tc.ISToUpdateBad, id[toUpdateI], nativeToUpdate)
		//			require.ErrStr(t, err)
		//		}
		//
		//		// update 1: ok
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.NoError(t, err, "what is an error on .Update()?")
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toUpdateI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		testData(t, keyFields, []string{id[toUpdateI]}, toUpdateResult, data, false, description, "on .Read() after Update()")
		//
		//		// update 2: ok
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.NoError(t, err, "what is an error on .Update()?")
		//		nativeToRead, err = tc.Read(tc.ISToUpdate, id[toUpdateI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		testData(t, keyFields, []string{id[toUpdateI]}, toUpdateResult, data, false, description, "on .Read() after Update()")
		//
		//		toUpdate := Tag{}
		//		for k, v := range toUpdateResult {
		//			toUpdate[k] = v
		//		}
		//
		//		// can't duplicate uniques fields
		//		for _, key := range uniquesUpdatable {
		//			toUpdate[key] = toSave[0][key]
		//			nativeToUpdate, err := tc.ItemToNative(toUpdate)
		//			require.NoError(t, err)
		//			err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//			require.ErrStr(t, err)
		//			toUpdate[key] = toUpdateResult[key]
		//		}
		//
		//		// update 3: ok
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.NoError(t, err, "what is the error on .Update()?")
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toUpdateI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		testData(t, keyFields, []string{id[toUpdateI]}, toUpdateResult, data, false, description, "on .Read() after Update()")
		//
		//		// can't update absent record
		//		toUpdate[keyFields[0]] += "123"
		//		nativeToUpdate, err = tc.ItemToNative(toUpdate)
		//		require.NoError(t, err)
		//		err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//		require.ErrStr(t, err)
		//	}
		//

		//toUpdateResult := tc.ToUpdate
		//for _, f := range description.FieldsArr {
		//	if !f.Creatable {
		//		toUpdateResult[f.Key] = data[f.Key]
		//	}
		//}

		// test ListTags -------------------------------------------------------------------------------------

		//if !tc.ExcludeListTest {
		//	var ids []common.Key
		//	for _, idi := range id {
		//		ids = append(ids, idi)
		//	}
		//
		//	if !tc.ExpectedReadOk {
		//		// TODO: selector.InStr(keyFields[0], ids...)
		//		briefsAll, err := tc.ListTags(nil, nil)
		//
		//		require.Equal(t, 0, len(briefsAll), "why len(dataAll) is not zero after .ListTags()?")
		//		require.ErrStr(t, err)
		//		continue
		//	}
		//
		//	// TODO: selector.InStr(keyFields[0], ids...)

		briefsAll, err := tc.List(nil, &crud.GetOptions{ActorKey: actorKey, OrderBy: []string{"id"}})
		require.NoError(t, err)
		require.True(t, len(briefsAll) == numRepeats)

		Compare(t, tc, &briefsAll[toReadI], toSave[i], l)
		Compare(t, tc, &briefsAll[toUpdateI], tc.ToUpdate, l)

		// test .Delete --------------------------------------------------------------------------------------

		err = tc.Remove(id[toDeleteI], &crud.RemoveOptions{ActorKey: actorKey})
		require.NoError(t, err)

		readDeleted, err := tc.Read(id[toDeleteI], &crud.GetOptions{ActorKey: actorKey})
		require.Error(t, err)
		require.Nil(t, readDeleted)

		briefsAll, err = tc.List(nil, &crud.GetOptions{ActorKey: actorKey})
		require.NoError(t, err)
		require.True(t, len(briefsAll) == numRepeats-1)

		//	if !tc.ExcludeRemoveTest {
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//		require.NoError(t, err, "what is the error on .Read() after Update()?")
		//		data, err = tc.NativeToItem(nativeToRead)
		//		require.NoError(t, err)
		//		require.Equal(t, id[toDeleteI], data[keyFields[0]])
		//
		//		if tc.ExpectedRemoveErr != nil {
		//			err = tc.Delete(tc.ISToDelete, id[toDeleteI])
		//			require.ErrStr(t, err, "where is an error on .DeleteList()?")
		//			nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//			require.NoError(t, err, "what is the error on .Read() after Update()?")
		//			data, err = tc.NativeToItem(nativeToRead)
		//			require.NoError(t, err)
		//			require.Equal(t, id[toDeleteI], data[keyFields[0]])
		//			continue
		//		}
		//
		//		if tc.ISToDeleteBad != nil {
		//			err = tc.Delete(*tc.ISToDeleteBad, id[toDeleteI])
		//			require.ErrStr(t, err, "where is an error on .DeleteList()?")
		//			nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//			require.NoError(t, err, "what is the error on .Read() after Update()?")
		//			data, err = tc.NativeToItem(nativeToRead)
		//			require.NoError(t, err)
		//			require.Equal(t, id[toDeleteI], data[keyFields[0]])
		//		}
		//
		//		err = tc.Delete(tc.ISToDelete, id[toDeleteI])
		//		require.NoError(t, err, "what is the error on .DeleteList()?")
		//
		//		nativeToRead, err = tc.Read(tc.ISToRead, id[toDeleteI])
		//
		//		// it depends on implementation
		//		// require.ErrStr(t, err, "where is an error on .Read() after DeleteList()?")
		//
		//		require.Nil(t, nativeToRead)
		//	}

	}
}
