import { createSlice } from '@reduxjs/toolkit';

interface WebCliState {
    isOpen: boolean
}
  

const initialState: WebCliState = {
    isOpen: false
};


export const webconsoleSlice = createSlice({
    name: 'webconsole',
    initialState,
    reducers: {
      toggleCli: (state) => {
        state.isOpen = !state.isOpen
      },
    },
  });

export const { toggleCli } = webconsoleSlice.actions;

export default webconsoleSlice.reducer;