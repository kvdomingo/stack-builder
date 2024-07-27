import { render } from "ink";
import meow from "meow";
import React from "react";
import App from "./app.js";

meow(
  `
	Usage
	  $ nodejs

	Options
		--name  Your name

	Examples
	  $ nodejs --name=Jane
	  Hello, Jane
`,
  {
    importMeta: import.meta,
  },
);

render(<App />);
