import React from 'react';
import { Calculator } from './components/Calculator';
import './styles/App.css';

export const App: React.FC = () => {
  return (
    <div className="app-layout">
      <Calculator />
    </div>
  );
};

export default App;