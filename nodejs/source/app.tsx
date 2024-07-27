import { flavors } from "@catppuccin/palette";
import { Box, Text } from "ink";
import SelectInput from "ink-select-input";
import React, { useState } from "react";
import { FRAMEWORKS } from "./constants.js";

const formKeyLabelMapping: { name: string; label: string }[] = [
  {
    name: "isTypeScript",
    label: "Use TypeScript",
  },
  {
    name: "framework",
    label: "Framework",
  },
  {
    name: "uiFramework",
    label: "UI Framework",
  },
  {
    name: "cssFramework",
    label: "CSS Framework",
  },
  {
    name: "database",
    label: "Database",
  },
  {
    name: "orm",
    label: "ORM",
  },
  {
    name: "cloudPlatform",
    label: "Cloud Platform",
  },
  {
    name: "isDocker",
    label: "Use Docker",
  },
];

const formKeys = formKeyLabelMapping.map(k => k.name);

interface Form {
  isTypeScript: boolean;
  framework: string;
  uiFramework: string;
  cssFramework: string;
  database: string;
  orm: string;
  cloudPlatform: string;
  isDocker: boolean;
}

const initialForm: Form = {
  isTypeScript: true,
  framework: "None",
  uiFramework: "None",
  cssFramework: "None",
  database: "None",
  orm: "None",
  cloudPlatform: "None",
  isDocker: false,
};

const palette = flavors.mocha.colors;

export default function App() {
  const [form, setForm] = useState({ ...initialForm });
  const [activeElement, setActiveElement] = useState(0);

  return (
    <Box width="100%">
      <Box flexDirection="column" gap={1} width="50%">
        <Box flexDirection="column">
          <Text backgroundColor={palette.mauve.hex} color={palette.crust.hex}>
            Do you want to use TypeScript?
          </Text>
          <SelectInput
            isFocused={formKeys[activeElement] === "isTypeScript"}
            items={[
              { label: "Yes", value: true },
              { label: "No", value: false },
            ]}
            onSelect={value => {
              setForm(form => ({ ...form, isTypeScript: value.value }));
              setActiveElement(prev => prev + 1);
            }}
          />
        </Box>

        <Box flexDirection="column">
          <Text backgroundColor={palette.mauve.hex} color={palette.crust.hex}>
            Choose a framework
          </Text>
          <SelectInput
            isFocused={formKeys[activeElement] === "framework"}
            items={FRAMEWORKS.map(fw => ({ label: fw, value: fw }))}
            onSelect={value => {
              setForm({ ...form, framework: value.value });
              setActiveElement(prev => prev + 1);
            }}
          />
        </Box>
      </Box>

      <Box
        width="50%"
        borderStyle="single"
        borderColor={palette.mauve.hex}
        flexDirection="column"
        gap={1}
      >
        <Text bold underline>
          Your Tech Stack
        </Text>

        {Object.entries(form).map(([key, value]) => (
          <Box flexDirection="column" key={key}>
            <Text color={palette.mauve.hex}>
              {formKeyLabelMapping.find(k => k.name === key)?.label ?? ""}
            </Text>
            <Text>{String(value)}</Text>
          </Box>
        ))}
      </Box>
    </Box>
  );
}
