import React from 'react';
// import { Counter } from './features/counter/Counter';
import { Navbar, Nav, Button } from 'react-bootstrap';

import { Database } from './features/database/Database';
// @ts-ignore  
import WebConsole from './features/web-console/WebConsole'
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import { useDispatch } from 'react-redux';
import { toggleCli } from './features/web-console/webconsoleSlice';

function App() {
  const dispatch = useDispatch();
  return (
    <div className="App">
      <Navbar bg="dark" variant="dark">
        <Navbar.Brand href="#home">Aoimdb</Navbar.Brand>
        <Nav className="mr-auto">
          <Nav.Item>        
            <Button onClick={() => dispatch(toggleCli())}> CLI </Button>
          </Nav.Item>
        </Nav>
        <Nav>
          <Nav.Link href="#features">Features</Nav.Link>
          <Nav.Link href="#about">About</Nav.Link>
        </Nav>
      </Navbar>
      <Navbar bg="dark" variant="dark" >
        <WebConsole />
      </Navbar>
      <header className="App-header">
        
        {/* <Counter /> */}
        <Database />
      </header>
    </div>
  );
}

export default App;
