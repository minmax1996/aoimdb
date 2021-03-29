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
      addToSet: (state, action: PayloadAction<string>) => {
        state.sets.push({key: action.payload, value: Object(action.payload)})
      },
    //   // Use the PayloadAction type to declare the contents of `action.payload`
    //   incrementByAmount: (state, action: PayloadAction<number>) => {
    //     state.value += action.payload;
    //   },
    },
  });

export const { addKey, addToSet } = databaseSlice.actions;

export const selectSetKeys = (state: RootState) => state.database.keys;

export default databaseSlice.reducer;