import React from 'react';
import { Counter } from './features/counter/Counter';
import { Database } from './features/database/Database';
import './App.css';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <Counter />
        <Database />
      </header>
    </div>
  );
}

export default App;
