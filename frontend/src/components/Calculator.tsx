import React, { useState } from 'react';
import { type Operation } from '../types/calculator';
import { calculatorClient } from '../api/calculatorClient';
import { Display } from './Display';
import { Keypad } from './Keypad';

export const Calculator: React.FC = () => {
  const [operandA, setOperandA] = useState<string>('');
  const [operandB, setOperandB] = useState<string>('');
  const [operation, setOperation] = useState<Operation>('add');
  const [result, setResult] = useState<number | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const validateInputs = (): string | null => {
    if (operandA.trim() === '') {
      return 'Please enter the first number.';
    }

    const numA = Number(operandA);
    if (isNaN(numA)) {
      return 'First value must be a valid number.';
    }

    if (operation !== 'sqrt') {
      if (operandB.trim() === '') {
        return 'Please enter the second number.';
      }

      const numB = Number(operandB);
      if (isNaN(numB)) {
        return 'Second value must be a valid number.';
      }

      if (operation === 'divide' && numB === 0) {
        return 'Cannot divide by zero.';
      }
    } else {
      if (numA < 0) {
        return 'Cannot calculate square root of a negative number.';
      }
    }

    return null;
  };

  const handleCalculate = async () => {
    setErrorMessage(null);
    setResult(null);

    const validationError = validateInputs();
    if (validationError) {
      setErrorMessage(validationError);
      return;
    }

    const numA = Number(operandA);
    const numB = operation !== 'sqrt' ? Number(operandB) : undefined;

    setIsLoading(true);

    try {
      const response = await calculatorClient.calculate(operation, numA, numB);

      if (response.ok) {
        setResult(response.value.result);
        setErrorMessage(null);
      } else {
        setErrorMessage(response.error.message || 'An error occurred.');
      }
    } catch {
      setErrorMessage('Unexpected error during calculation.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleClear = () => {
    setOperandA('');
    setOperandB('');
    setResult(null);
    setErrorMessage(null);
  };

  return (
    <main className="calculator-container" aria-label="REST Calculator">
      <header className="calculator-header">
        <h1>REST Calculator</h1>
      </header>

      <Display
        operandA={operandA}
        operandB={operandB}
        operation={operation}
        result={result}
        errorMessage={errorMessage}
        isLoading={isLoading}
      />

      <Keypad
        operandA={operandA}
        operandB={operandB}
        operation={operation}
        onOperandAChange={(val) => {
          setOperandA(val);
          if (errorMessage) setErrorMessage(null);
        }}
        onOperandBChange={(val) => {
          setOperandB(val);
          if (errorMessage) setErrorMessage(null);
        }}
        onOperationChange={(op) => {
          setOperation(op);
          if (errorMessage) setErrorMessage(null);
        }}
        onCalculate={handleCalculate}
        onClear={handleClear}
        isLoading={isLoading}
      />
    </main>
  );
};