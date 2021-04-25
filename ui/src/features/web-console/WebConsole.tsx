import React, { useState } from "react";
import { Collapse } from "react-bootstrap";

import ReactConsole from "@webscopeio/react-console";
import { useDispatch, useSelector } from "react-redux";
import { addKey, clearKeys } from "../data-visualizator/databaseSlice";
import { selectIsOpen } from "../../app/store";

interface Command {
  description: string;
  fn: (...args: any[]) => Promise<any>;
}

export default function WebConsole() {
  const [history, setHistory] = useState<string[]>([]);
  const isOpens = useSelector(selectIsOpen);
  const dispatch = useDispatch();
  return (
    <Collapse in={isOpens}>
      <div id="br-dark" style={{ width: "100%" }}>
        <ReactConsole
          autoFocus
          prompt=">>>"
          welcomeMessage="Welcome to Cli
          available commands: 
          echo args... ; keys ; clearkeys ; history
          "
          // wrapperClassName="bg-dark"
          wrapperStyle={{ textAlign: "left", marginBottom: "10px" }}
          history={history}
          onAddHistoryItem={(item) => setHistory([...history, item])}
          commands={{
            echo: {
              fn: (...args: any[]) => {
                return new Promise<any>((resolve, _) => {
                  resolve(`${args.join(" ")}`);
                  dispatch(addKey(`${args.join(" ")}`));
                });
              },
            },
            clearkeys: {
              fn: () => {
                return new Promise<void>((resolve) => {
                  dispatch(clearKeys());
                  resolve();
                });
              },
            },
            keys: {
              fn: () => {
                return new Promise<void>((resolve) => {
                  dispatch(clearKeys());
                  dispatch(addKey("get keys from database"));
                  resolve();
                });
              },
            },
            history: {
              description: "History",
              fn: () =>
                new Promise((resolve) => {
                  resolve(`${history.join("\r\n")}`);
                }),
            },
          }}
        />
      </div>
    </Collapse>
  );
}
