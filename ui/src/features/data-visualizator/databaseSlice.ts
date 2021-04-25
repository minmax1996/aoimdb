import { createSlice, PayloadAction } from "@reduxjs/toolkit";
import { AppThunk, RootState } from "../../app/store";

interface Sets {
  key: string;
  value: object;
}

interface Tables {
  columnNames: string[];
  datas: string[][];
}

interface DatabaseState {
  keys: string[];
  sets: Sets[];
  tables: Tables[];
}

const initialState: DatabaseState = {
  keys: ["key1", "key2", "key3"],
  sets: [{ key: "set1", value: Object("set1val1") }],
  tables: [
    {
      columnNames: ["id", "col1", "col2"],
      datas: [
        ["1", "col1val", "col2val"],
        ["2", "col1val2", "col2val2"],
      ],
    },
  ],
};

export const databaseSlice = createSlice({
  name: "database",
  initialState,
  reducers: {
    addKey: (state, action: PayloadAction<string>) => {
      state.keys.push(action.payload);
    },
    clearKeys: (state, action: PayloadAction<void>) => {
      state.keys = [];
    },
    addToSet: (state, action: PayloadAction<Sets>) => {
      state.sets.push(action.payload);
    },
  },
});

export const { addKey, addToSet, clearKeys } = databaseSlice.actions;

export const selectSetKeys = (state: RootState) => state.database.keys;

export default databaseSlice.reducer;
