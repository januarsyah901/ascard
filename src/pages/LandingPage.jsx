import { Link } from 'react-router-dom';
import Navigation from '../components/Navigation';
import Footer from '../components/Footer';
import { Gem, Guardian } from '../components/OlympusArtwork';
import { useLanguage } from '../context/LanguageContext';
import { useCurrency } from '../context/CurrencyContext';
import '../App.css';

const symbols = [0,6,7,5,1,2,1,3,3,3,1,4,5,4,0,5,5,2,1,4,1,7,7,4,5,6,2,0,5,2];

export default function LandingPage() {
  const { lang } = useLanguage();
  const { formatValue } = useCurrency();
  const say = (id, en) => lang === 'en' ? en : id;
  return <div className="bi-lobby ol-lobby"><Navigation/><main id="main-content" className="pt-16">
    <div className="ol-education-strip">◇ {say('BONGKAR ILUSI — PLATFORM EDUKASI ANTI-JUDI', 'BONGKAR ILUSI — ANTI-GAMBLING EDUCATION')} <span>• {say('TANPA DEPOSIT · TANPA UANG ASLI', 'NO DEPOSITS · NO REAL MONEY')}</span></div>
    <section className="ol-temple" aria-labelledby="ol-heading"><div className="ol-clouds"/><div className="ol-column ol-column-left"/><div className="ol-column ol-column-right"/>
      <div className="ol-stage-heading"><span>✦ {say('DI BALIK GEMERLAP, ADA MATEMATIKA', 'BEHIND THE GLAMOUR LIES MATHEMATICS')} ✦</span><h1 id="ol-heading">{say('BONGKAR ILUSI', 'BONGKAR ILUSI')}</h1><p>{say('Tampilan yang menggoda. Kenyataan yang berbeda.', 'An alluring appearance. A different reality.')}</p></div>
      <div className="ol-stage">
        <aside className="ol-left-panels"><Link to="/simulation" className="ol-feature-panel"><span className="ol-mini-label">PROBABILITY LAB</span><strong>{say('COBA', 'TRY')}<br/>{say('SIMULASI', 'SIMULATION')}</strong><small>{say('100% VIRTUAL', '100% VIRTUAL')}</small><span className="ol-panel-action">{say('MASUK LAB', 'ENTER LAB')} →</span></Link><Link to="/knowledge-base" className="ol-feature-panel ol-panel-pink"><span className="ol-mini-label">KNOWLEDGE BASE</span><strong>{say('BONGKAR', 'UNCOVER')}<br/>{say('MITOSNYA', 'THE MYTHS')}</strong><small>{say('BACA FAKTANYA', 'READ THE FACTS')}</small><span className="ol-panel-action">{say('BUKA ARTIKEL', 'READ ARTICLES')} →</span></Link><div className="ol-side-note"><span>ϟ</span><b>{say('KENALI ILUSINYA', 'KNOW THE ILLUSION')}</b><p>{say('Simbol di layar ini dekoratif. Bukan hasil putaran dan bukan hadiah.', 'These symbols are decorative. They are not spin results or prizes.')}</p></div></aside>
        <div className="ol-game-frame"><div className="ol-corner ol-corner-a">◆</div><div className="ol-corner ol-corner-b">◆</div><div className="ol-board-caption">{say('LAYAR ILUSTRASI · BUKAN PERMAINAN AKTIF', 'ILLUSTRATIVE DISPLAY · NOT AN ACTIVE GAME')}</div><div className="ol-demo-grid" aria-hidden="true">{symbols.map((kind, i) => <div className="ol-symbol-cell" key={i}><Gem kind={kind}/></div>)}</div><div className="ol-board-bottom">✦ {say('PAHAMI PELUANG. KENALI JEBAKANNYA.', 'UNDERSTAND THE ODDS. RECOGNIZE THE TRAP.')} ✦</div></div>
        <aside className="ol-character-panel"><div className="ol-game-title"><span>GATES OF</span><strong>ILLUSION</strong><b>EDUCATION</b></div><Guardian/><div className="ol-character-caption">{say('GEMERLAPNYA NYATA.', 'THE GLAMOUR IS REAL.')}<br/><strong>{say('JANJINYA BELUM TENTU.', 'THE PROMISE MAY NOT BE.')}</strong></div></aside>
      </div>
      <div className="ol-control-bar"><div><small>{say('KREDIT ILUSTRASI', 'ILLUSTRATIVE CREDIT')}</small><strong>{formatValue(333.333333)}</strong><span>{say('Saldo virtual • Tidak dapat diuangkan', 'Virtual balance • Cannot be cashed out')}</span></div><div className="ol-control-message"><strong>{say('BUKAN SOAL HOKI.', 'IT’S NOT ABOUT LUCK.')}</strong><span>{say('Yuk, lihat cara sistemnya bekerja.', 'See how the system works.')}</span></div><Link to="/simulation" className="ol-start-control" aria-label={say('Buka simulator edukasi', 'Open educational simulator')}><span aria-hidden="true">⟳</span><b>{say('MULAI SIMULASI', 'START SIMULATION')}</b></Link></div>
    </section>
    <div className="ol-bottom-content"><section className="ol-reality-banner"><div><span className="ol-mini-label">REALITY CHECK</span><h2>{say('Terlihat seperti slot. Dibuat untuk membuka mata.', 'Looks like a slot. Built to open your eyes.')}</h2><p>{say('Tidak ada deposit, hadiah uang, atau janji kemenangan. Gunakan simulator untuk mengamati saldo virtual, lalu pelajari probabilitas dan bias yang membentuk keputusan.', 'No deposits, cash prizes, or winning promises. Use the simulator to observe a virtual balance, then explore probability and the biases behind decisions.')}</p></div><Link to="/intervention" className="ol-gold-button">{say('AMBIL JEDA & REFLEKSI', 'TAKE A BREAK & REFLECT')} →</Link></section><div className="ol-learning-grid">{[['01','/simulation',say('Masuk ke laboratorium', 'Enter the laboratory'),say('Amati hasil putaran dan riwayat saldo virtual.', 'Observe spin results and virtual balance history.')],['02','/knowledge-base',say('Baca di balik simbol', 'Read behind the symbols'),say('Pelajari matematika dan psikologi perjudian.', 'Explore the mathematics and psychology of gambling.')],['03','/intervention',say('Putus rantai ilusinya', 'Break the illusion'),say('Kenali mitos dan ambil ruang untuk berpikir.', 'Recognize myths and make room to reflect.')]].map(([n,path,title,text])=><Link to={path} className="ol-learning-card" key={n}><span>{n}</span><h3>{title} ↗</h3><p>{text}</p></Link>)}</div></div>
  </main><Footer/></div>;
}
