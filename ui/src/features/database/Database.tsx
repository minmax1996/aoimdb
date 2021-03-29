import React, { useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  addKey,
  addToSet,
  selectSetKeys,
} from './databaseSlice';
import styles from '../counter/Counter.module.css';



export function Database() {
    const count = useSelector(selectSetKeys);
    const dispatch = useDispatch();
    const [keyAmount, setKeyAmount] = useState('');
    const [keySetAmount, setSetKeyAmount] = useState('');
  
    return (
      <div>
        <div className={styles.row}>
        <input
          className={styles.textbox}
          aria-label="Set key"
          value = {keyAmount}
          onChange={e => setKeyAmount(e.target.value)}
        />
        <input
          className={styles.textbox}
          aria-label="Set Value key"
          value = {keySetAmount}
          onChange={e => setSetKeyAmount(e.target.value)}
        />
          <button
            className={styles.button}
            aria-label="Decrement value"
            onClick={() => dispatch(addKey(keyAmount))}
          >
            1
          </button>
          <button
            className={styles.button}
            aria-label="Decrement value"
            onClick={() => dispatch(addKey(keySetAmount))}
          >
            2
          </button>
          <span className={styles.value}>{count}</span>
        </div>
      </div>
    );
  }
  