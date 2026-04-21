import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import KnowledgeBase from './pages/KnowledgeBase';
import SimulationLab from './pages/SimulationLab';
import Intervention from './pages/Intervention';

import { LanguageProvider } from './context/LanguageContext';
import { CurrencyProvider } from './context/CurrencyContext';

function App() {
  return (
    <LanguageProvider>
      <CurrencyProvider>
        <Router>
          <Routes>
            <Route path="/" element={<LandingPage />} />
            <Route path="/knowledge-base" element={<KnowledgeBase />} />
            <Route path="/simulation" element={<SimulationLab />} />
            <Route path="/intervention" element={<Intervention />} />
          </Routes>
        </Router>
      </CurrencyProvider>
    </LanguageProvider>
  );
}

export default App;
