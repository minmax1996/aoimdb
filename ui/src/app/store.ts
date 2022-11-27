import { configureStore, ThunkAction, Action } from "@reduxjs/toolkit";
import databaseReducer from "../features/data-visualizator/databaseSlice";
import webconsoleReduser from "../features/web-console/webconsoleSlice";

export const store = configureStore({
  reducer: {
    database: databaseReducer,
    cli: webconsoleReduser,
  },
});

export const selectIsOpen = (state: RootState) => state.cli.isOpen;

export type RootState = ReturnType<typeof store.getState>;
export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  RootState,
  unknown,
  Action<string>
>;
