import React, { useState } from "react";
import { Tab, Table, Tabs } from "react-bootstrap";
import { useSelector, useDispatch } from "react-redux";
import {
  selectSetKeys,
  selectSets,
  selectTables,
  Tables,
} from "./databaseSlice";

export function DataVisualizator() {
  const keys = useSelector(selectSetKeys);
  const sets = useSelector(selectSets);
  return (
    <div>
      <p>Visualization</p>
      <Tabs defaultActiveKey="table" id="uncontrolled-tab-example">
        <Tab eventKey="keys" title="Keys">
          <p>Keys</p>
        </Tab>
        <Tab eventKey="table" title="Tables">
          {TableVisualizator(useSelector(selectTables))}
        </Tab>
      </Tabs>
    </div>
  );
}

export function TableVisualizator(params: Tables[]) {
  return params.map((table) => (
    <Table striped bordered hover>
      <thead>
        {table.columnNames.map((col: string) => (
          <th>{col}</th>
        ))}
      </thead>
      <tbody>
        {table.datas.map((row: string[]) => (
          <tr>
            {row.map((col: string) => (
              <td>{col}</td>
            ))}
          </tr>
        ))}
      </tbody>
    </Table>
  ));
}
