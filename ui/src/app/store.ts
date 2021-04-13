import { configureStore, ThunkAction, Action } from '@reduxjs/toolkit';
// @ts-ignore
import swaggerClient from 'redux-swagger-client'
import thunk from 'redux-thunk';
import counterReducer from '../features/counter/counterSlice';
import databaseReducer from '../features/database/databaseSlice';

export const store = configureStore({
  reducer: {
    counter: counterReducer,
    database: databaseReducer,
  },
  middleware: (getDefaultMiddleware) => getDefaultMiddleware().concat(thunk).concat(swaggerClient({url:'0.0.0.0:3000/static/command.swagger.json'}))
});

export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
