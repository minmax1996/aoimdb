import React, { useState } from "react";
import { useSelector, useDispatch } from "react-redux";
import { addKey, addToSet, selectSetKeys } from "./databaseSlice";
import styles from "../counter/Counter.module.css";

export function DataVisualizator() {
  const keys = useSelector(selectSetKeys);
  //sconst dispatch = useDispatch();
  // const [keyAmount, setKeyAmount] = useState('');
  // const [keySetAmount, setSetKeyAmount] = useState('');
  return (
    <div>
      <div className={styles.row}>
        <p>
          {" "}
          <p style={{ fontSize: 22 }}>KEYS: </p>
          {keys.map((key) => {
            return (
              <div>
                {key}
                <br />
              </div>
            );
          })}
        </p>
      </div>
    </div>
  );
}
