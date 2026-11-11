import { Link } from 'react-router-dom';
import { useLanguage } from '../context/LanguageContext';
import '../App.css';

export default function Footer() {
  const { lang, t } = useLanguage();
  const en = lang === 'en';
  return (
    <footer className="bi-footer">
      <div className="bi-container py-10">
        <div className="flex flex-col md:flex-row justify-between gap-8"><div className="max-w-md"><Link to="/" className="bi-brand">BONGKAR <b>ILUSI</b></Link><p className="mt-4 text-sm leading-relaxed text-stone-400">{en ? 'An educational platform exploring probability and the psychology of gambling. Virtual simulations with no deposits or cash prizes.' : 'Platform edukasi untuk memahami probabilitas dan psikologi perjudian. Simulasi virtual tanpa deposit dan hadiah uang.'}</p></div><nav aria-label={en ? 'Footer navigation' : 'Navigasi footer'} className="flex flex-wrap items-start gap-x-6 gap-y-3 text-sm text-stone-300"><Link to="/simulation">{t.nav.simulator}</Link><Link to="/knowledge-base">{t.nav.articles}</Link><Link to="/intervention">{en ? 'Reflection & myths' : 'Refleksi & mitos'}</Link></nav></div>
        <div className="bi-footer-bottom flex flex-col sm:flex-row justify-between gap-3 mt-8 pt-5 text-xs text-stone-400"><span>© 2026 Bongkar Ilusi.</span><span>{en ? 'EDUCATIONAL DEMO • NO GAMBLING LICENSE CLAIM' : 'DEMO EDUKATIF • TANPA KLAIM LISENSI JUDI'}</span></div>
      </div>
    </footer>
  );
}
