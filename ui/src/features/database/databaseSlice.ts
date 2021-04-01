import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { AppThunk, RootState } from '../../app/store';

interface Sets {
    key: string,
    value: object,
}

interface DatabaseState {
  keys: string[];
  sets: Sets[];
}

const initialState: DatabaseState = {
    keys: [],
    sets: [],
};

export const databaseSlice = createSlice({
    name: 'database',
    initialState,
    reducers: {
      addKey: (state, action: PayloadAction<string>) => {
        state.keys.push(action.payload)
      },
      clearKeys: (state, action: PayloadAction<void>) => {
        state.keys=[]
      },
      addToSet: (state, action: PayloadAction<Sets>) => {
        state.sets.push(action.payload)
      },
    },
  });

export const { addKey, addToSet, clearKeys} = databaseSlice.actions;

export const selectSetKeys = (state: RootState) => state.database.keys;

export default databaseSlice.reducer;