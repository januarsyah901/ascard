import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import LandingPage from './pages/LandingPage';
import KnowledgeBase from './pages/KnowledgeBase';
import SimulationLab from './pages/SimulationLab';
import Intervention from './pages/Intervention';
import AdminLogin from './pages/AdminLogin';
import AdminDashboard from './pages/AdminDashboard';
import ArticleDetail from './pages/ArticleDetail';

import { LanguageProvider } from './context/LanguageContext';
import { CurrencyProvider } from './context/CurrencyContext';
import { AuthProvider } from './context/AuthContext';

function App() {
  return (
    <LanguageProvider>
      <CurrencyProvider>
        <AuthProvider>
          <Router>
            <Routes>
              <Route path="/" element={<LandingPage />} />
              <Route path="/knowledge-base" element={<KnowledgeBase />} />
              <Route path="/simulation" element={<SimulationLab />} />
              <Route path="/intervention" element={<Intervention />} />
              <Route path="/admin/login" element={<AdminLogin />} />
              <Route path="/admin/dashboard" element={<AdminDashboard />} />
              <Route path="/article/:id" element={<ArticleDetail />} />
            </Routes>
          </Router>
        </AuthProvider>
      </CurrencyProvider>
    </LanguageProvider>
  );
}

export default App;
