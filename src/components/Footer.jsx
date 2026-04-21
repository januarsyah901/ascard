import React from 'react';
import { Link } from 'react-router-dom';

export default function Footer() {
  return (
    <footer className="w-full py-12 px-6 border-t-0 bg-black text-white">
      <div className="flex flex-col md:flex-row justify-between items-center max-w-7xl mx-auto gap-8">
        <div className="font-black text-white text-xl tracking-tighter uppercase font-headline">
          BONGKAR ILUSI
        </div>
        <div className="flex flex-wrap justify-center gap-8">
          <Link to="/knowledge-base" className="text-neutral-400 hover:text-white transition-colors font-light tracking-widest text-[10px] uppercase font-headline">Terms of Inquiry</Link>
          <Link to="/knowledge-base" className="text-neutral-400 hover:text-white transition-colors font-light tracking-widest text-[10px] uppercase font-headline">Methodology</Link>
          <Link to="/knowledge-base" className="text-neutral-400 hover:text-white transition-colors font-light tracking-widest text-[10px] uppercase font-headline">Privacy Protocol</Link>
          <Link to="/intervention" className="text-neutral-400 hover:text-white transition-colors font-light tracking-widest text-[10px] uppercase font-headline">Contact Lab</Link>
        </div>
        <div className="font-light tracking-widest text-[10px] uppercase font-headline opacity-60 text-center md:text-right">
          © 2026 CLINICAL AUTHORITY. ALL RIGHTS RESERVED. EMPIRICAL MONOLITH DESIGN SYSTEM.
        </div>
      </div>
    </footer>
  );
}
