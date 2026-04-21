import React from 'react';
import { Link, NavLink } from 'react-router-dom';
import { useLanguage } from '../context/LanguageContext';
import { useCurrency } from '../context/CurrencyContext';

export default function Navigation() {
  const { lang, t, toggleLanguage } = useLanguage();
  const { currency, toggleCurrency } = useCurrency();
  
  const navLinkStyle = ({ isActive }) => 
    `transition-all duration-200 font-bold tracking-tighter uppercase font-headline text-sm pb-1 ${
      isActive ? 'text-black border-b-2 border-black' : 'text-neutral-500 hover:text-black'
    }`;

  const ToggleSwitch = ({ value, options, onChange, width = "w-20" }) => {
    const activeIndex = options.indexOf(value);
    return (
      <div 
        onClick={onChange}
        className={`relative flex items-center bg-neutral-100 p-1 cursor-pointer ${width} h-8 select-none border border-black/5 hover:border-black/20 transition-colors`}
      >
        <div 
          className="absolute h-6 bg-black transition-all duration-300 ease-out shadow-sm"
          style={{ 
            width: 'calc(50% - 4px)',
            left: activeIndex === 0 ? '4px' : 'calc(50%)'
          }}
        />
        {options.map((opt, i) => (
          <div 
            key={opt}
            className={`relative z-10 flex-1 text-center text-[9px] font-black tracking-widest uppercase transition-colors duration-300 ${
              activeIndex === i ? 'text-white' : 'text-neutral-400'
            }`}
          >
            {opt}
          </div>
        ))}
      </div>
    );
  };

  return (
    <nav className="fixed top-0 w-full z-50 h-16 bg-white/80 backdrop-blur-xl border-b border-black/5">
      <div className="flex justify-between items-center px-6 md:px-12 w-full mx-auto h-full">
        <div className="text-xl font-black tracking-tighter text-black uppercase font-headline">
          <Link to="/">BONGKAR ILUSI</Link>
        </div>
        <div className="hidden md:flex items-center gap-8">
          <NavLink to="/simulation" className={navLinkStyle}>{t.nav.simulator}</NavLink>
          <NavLink to="/knowledge-base" className={navLinkStyle}>{t.nav.library}</NavLink>
          <NavLink to="/intervention" className={navLinkStyle}>{t.nav.logic}</NavLink>
        </div>
        <div className="flex items-center gap-4">
          <Link to="/intervention">
            <button className="px-6 py-2 bg-primary text-on-primary font-bold tracking-tighter uppercase font-headline text-xs hover:bg-primary-fixed-dim transition-all">
              {t.nav.getHelp}
            </button>
          </Link>
          <div className="flex items-center gap-2 ml-2">
            <ToggleSwitch 
              value={lang.toUpperCase()} 
              options={['ID', 'EN']} 
              onChange={toggleLanguage} 
              width="w-24"
            />
            <ToggleSwitch 
              value={currency} 
              options={['IDR', 'USD']} 
              onChange={toggleCurrency} 
              width="w-24"
            />
          </div>
        </div>
      </div>
    </nav>
  );
}
