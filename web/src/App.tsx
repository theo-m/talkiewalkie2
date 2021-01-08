import React from "react";

import { ABClient } from "./TwServiceClientPb";
import { Person, AddressBook } from "./tw_pb";

import "./App.css";

const c = new ABClient("http://localhost:8080");

function App() {
  return (
    <div className="App">
      <button
        onClick={() =>
          c.get(new AddressBook(), { "custom-header-1": "hellooo" })
        }
      >
        ping
      </button>
    </div>
  );
}

export default App;
