import React from 'react';
import { type Operation } from '../types/calculator';

interface KeypadProps {
  operandA: string;
  operandB: string;
  operation: Operation;
  onOperandAChange: (val: string) => void;
  onOperandBChange: (val: string) => void;
  onOperationChange: (op: Operation) => void;
  onCalculate: () => void;
  onClear: () => void;
  isLoading: boolean;
}

const OPERATIONS: { id: Operation; label: string; symbol: string }[] = [
  { id: 'add', label: 'Add', symbol: '+' },
  { id: 'subtract', label: 'Subtract', symbol: '−' },
  { id: 'multiply', label: 'Multiply', symbol: '×' },
  { id: 'divide', label: 'Divide', symbol: '÷' },
  { id: 'power', label: 'Power', symbol: 'xʸ' },
  { id: 'percentage', label: 'Percentage', symbol: '%' },
  { id: 'sqrt', label: 'Square Root', symbol: '√x' },
];

export const Keypad: React.FC<KeypadProps> = ({
  operandA,
  operandB,
  operation,
  onOperandAChange,
  onOperandBChange,
  onOperationChange,
  onCalculate,
  onClear,
  isLoading,
}) => {
  const isUnary = operation === 'sqrt';

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      onCalculate();
    }
  };

  return (
    <div className="calculator-keypad" onKeyDown={handleKeyDown}>
      <div className="input-group">
        <label htmlFor="operand-a-input" className="input-label">
          {isUnary ? 'Number (x)' : 'First Number (a)'}
        </label>
        <input
          id="operand-a-input"
          type="number"
          step="any"
          className="number-input"
          value={operandA}
          onChange={(e) => onOperandAChange(e.target.value)}
          placeholder="0"
          disabled={isLoading}
          autoFocus
        />
      </div>

      <div className="operations-grid" role="group" aria-label="Select operation">
        {OPERATIONS.map((op) => (
          <button
            key={op.id}
            type="button"
            className={`operation-btn ${operation === op.id ? 'active' : ''}`}
            onClick={() => onOperationChange(op.id)}
            disabled={isLoading}
            aria-label={op.label}
            aria-pressed={operation === op.id}
          >
            {op.symbol}
          </button>
        ))}
      </div>

      {!isUnary && (
        <div className="input-group">
          <label htmlFor="operand-b-input" className="input-label">
            Second Number (b)
          </label>
          <input
            id="operand-b-input"
            type="number"
            step="any"
            className="number-input"
            value={operandB}
            onChange={(e) => onOperandBChange(e.target.value)}
            placeholder="0"
            disabled={isLoading}
          />
        </div>
      )}

      <div className="action-buttons">
        <button
          type="button"
          className="action-btn clear-btn"
          onClick={onClear}
          disabled={isLoading}
          aria-label="Clear inputs"
        >
          Clear
        </button>

        <button
          type="button"
          className="action-btn calculate-btn"
          onClick={onCalculate}
          disabled={isLoading}
          aria-label="Calculate result"
        >
          {isLoading ? 'Calculating...' : '='}
        </button>
      </div>
    </div>
  );
};