import { useState } from 'react';
import { Link, NavLink, useLocation } from 'react-router-dom';
import { useLanguage } from '../context/LanguageContext';
import { useCurrency } from '../context/CurrencyContext';
import '../App.css';

export default function Navigation() {
  const { lang, t, toggleLanguage } = useLanguage();
  const { currency, toggleCurrency } = useCurrency();
  const { pathname } = useLocation();
  const [open, setOpen] = useState(false);
  const links = [['/', t.nav.home], ['/simulation', t.nav.simulator], ['/knowledge-base', t.nav.articles], ['/intervention', t.nav.logic]];
  const renderLinks = () => links.map(([path, label]) => <NavLink key={path} to={path} end={path === '/'} onClick={() => setOpen(false)} className={({ isActive }) => `bi-nav-link ${isActive ? 'bi-nav-active' : ''}`}>{label}</NavLink>);
  return (
    <>
      {pathname === '/' && <a href="#main-content" className="bi-skip">{lang === 'en' ? 'Skip to content' : 'Langsung ke konten'}</a>}
      <header className="bi-navigation fixed top-0 inset-x-0 z-50" onKeyDown={event => { if (event.key === 'Escape') { setOpen(false); event.currentTarget.querySelector('[aria-controls="bi-mobile-menu"]')?.focus(); } }}>
        <nav aria-label={lang === 'en' ? 'Main navigation' : 'Navigasi utama'} className="bi-nav-inner">
          <Link to="/" className="bi-brand" onClick={() => setOpen(false)}><span className="bi-brand-mark" aria-hidden="true">B<span>i</span></span><span>BONGKAR <b>ILUSI</b><small>{lang === 'en' ? 'ANTI-GAMBLING EDUCATION' : 'EDUKASI ANTI-JUDI'}</small></span></Link>
          <div className="hidden lg:flex items-center gap-5">{renderLinks()}</div>
          <div className="flex items-center gap-2">
            <button type="button" onClick={toggleLanguage} className="bi-toggle" aria-label={lang === 'en' ? 'Switch to Indonesian' : 'Ganti ke bahasa Inggris'}>{lang.toUpperCase()} <span aria-hidden="true">⌄</span></button>
            <button type="button" onClick={toggleCurrency} className="bi-toggle" aria-label={lang === 'en' ? `Currency ${currency}; switch currency` : `Mata uang ${currency}; ganti mata uang`}>{currency}</button>
            <Link to="/intervention" className="bi-button bi-button-red hidden sm:inline-flex">{t.nav.getHelp}</Link>
            <button type="button" className="bi-toggle lg:hidden" aria-expanded={open} aria-controls="bi-mobile-menu" aria-label={lang === 'en' ? 'Toggle navigation' : 'Buka atau tutup navigasi'} onClick={() => setOpen(value => !value)}>{open ? '✕' : '☰'}</button>
          </div>
        </nav>
        {open && <div id="bi-mobile-menu" className="bi-mobile-menu lg:hidden">{renderLinks()}<Link to="/intervention" onClick={() => setOpen(false)} className="bi-button bi-button-red">{t.nav.getHelp}</Link></div>}
      </header>
    </>
  );
}
