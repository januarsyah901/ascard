import React from 'react';
import { Link } from 'react-router-dom';
import Navigation from '../components/Navigation';
import Footer from '../components/Footer';
import { useLanguage } from '../context/LanguageContext';

export default function KnowledgeBase() {
  const { t } = useLanguage();
  const [dynamicArticles, setDynamicArticles] = React.useState([]);

  React.useEffect(() => {
    fetch('http://localhost:5001/api/articles')
      .then(res => res.json())
      .then(data => {
        setDynamicArticles(data.filter(a => a.published));
      })
      .catch(err => console.error('Error fetching articles:', err));
  }, []);

  return (
    <div className="min-h-screen flex flex-col bg-surface text-on-surface">
      <Navigation />
      <main className="flex-grow pt-24 pb-16 px-6 md:px-12 max-w-7xl mx-auto w-full">
        {/* Hero Header */}
        <header className="mb-16">
          <div className="inline-block bg-primary text-on-primary-fixed px-3 py-1 mb-6 font-label text-[10px] font-bold tracking-widest uppercase">
            {t.knowledge.subtitle}
          </div>
          <h1 className="font-headline font-black text-6xl md:text-8xl tracking-tighter uppercase leading-none mb-8">
            {t.knowledge.title.split(' ')[0]} <br /> {t.knowledge.title.split(' ')[1]}
          </h1>
          <div className="editorial-grid gap-6 grid grid-cols-12">
            <p className="col-span-12 md:col-span-6 font-body text-xl text-on-surface-variant leading-relaxed">
              {t.knowledge.desc}
            </p>
            <div className="col-span-12 md:col-span-6 flex flex-wrap gap-2 items-end justify-start md:justify-end">
              <span className="px-4 py-2 bg-surface-container-highest font-label font-bold text-[11px] uppercase tracking-widest">{t.knowledge.psychology}</span>
              <span className="px-4 py-2 bg-surface-container font-label font-medium text-[11px] uppercase tracking-widest text-on-surface-variant">{t.knowledge.math}</span>
              <span className="px-4 py-2 bg-surface-container font-label font-medium text-[11px] uppercase tracking-widest text-on-surface-variant">{t.knowledge.finance}</span>
            </div>
          </div>
        </header>
        
        {/* Feature Bento Grid */}
        <div className="grid grid-cols-12 gap-1 mb-1 bg-outline-variant/20">
          {/* Main Featured Article */}
          <article className="col-span-12 md:col-span-8 bg-surface-container-lowest p-8 md:p-12 relative group overflow-hidden">
            <div className="relative z-10">
              <div className="flex items-center gap-4 mb-6">
                <span className="font-label font-black text-[11px] uppercase tracking-[0.2em] text-primary">{t.knowledge.featured}</span>
                <div className="h-[1px] flex-grow bg-outline-variant/30"></div>
                <span className="font-label font-bold text-[11px] uppercase tracking-widest flex items-center gap-2">
                  <span className="material-symbols-outlined text-sm">schedule</span> {t.knowledge.readTime}
                </span>
              </div>
              <h2 className="font-headline font-extrabold text-4xl md:text-5xl tracking-tighter uppercase mb-6 leading-tight group-hover:text-primary-fixed-dim transition-colors text-on-surface">
                {t.knowledge.anatomyTitle}
              </h2>
              <p className="font-body text-lg text-on-surface-variant mb-8 max-w-2xl">
                {t.knowledge.anatomyDesc}
              </p>
              <div className="h-64 w-full bg-surface-container mb-8 flex items-center justify-center relative">
                <div className="absolute inset-0 opacity-10 bg-[radial-gradient(circle_at_center,_var(--tw-gradient-stops))] from-primary to-transparent"></div>
                <div className="w-full h-full flex items-end justify-between px-12 pb-8">
                  <div className="w-8 bg-primary" style={{ height: '90%' }}></div>
                  <div className="w-8 bg-neutral-400" style={{ height: '70%' }}></div>
                  <div className="w-8 bg-neutral-300" style={{ height: '55%' }}></div>
                  <div className="w-8 bg-neutral-200" style={{ height: '40%' }}></div>
                  <div className="w-8 bg-neutral-100" style={{ height: '25%' }}></div>
                  <div className="w-8 bg-neutral-50" style={{ height: '10%' }}></div>
                </div>
                <div className="absolute top-4 left-6 font-label font-bold text-[10px] uppercase text-neutral-400">Model Probabilitas: Probabilitas Kebangkrutan vs Waktu</div>
              </div>
              <button className="font-label font-black text-xs uppercase tracking-widest border-b-2 border-primary pb-1">{t.knowledge.readMore}</button>
            </div>
          </article>
          
          {/* Secondary Stat Block */}
          <aside className="col-span-12 md:col-span-4 bg-primary text-on-primary-fixed p-8 flex flex-col justify-between">
            <div>
              <h3 className="font-headline font-black text-2xl tracking-tighter uppercase mb-4">{t.knowledge.realityStats}</h3>
              <div className="space-y-8 mt-8">
                <div>
                  <span className="block font-headline font-extrabold text-5xl tracking-tighter">99.8%</span>
                  <span className="block font-label text-[10px] uppercase tracking-widest mt-2 opacity-70">{t.knowledge.stat1}</span>
                </div>
                <div>
                  <span className="block font-headline font-extrabold text-5xl tracking-tighter">0.000001%</span>
                  <span className="block font-label text-[10px] uppercase tracking-widest mt-2 opacity-70">{t.knowledge.stat2}</span>
                </div>
                <div>
                  <span className="block font-headline font-extrabold text-5xl tracking-tighter">1.4ms</span>
                  <span className="block font-label text-[10px] uppercase tracking-widest mt-2 opacity-70">{t.knowledge.stat3}</span>
                </div>
              </div>
            </div>
            <div className="pt-8 border-t border-white/10">
              <p className="font-label text-[9px] uppercase tracking-widest opacity-50">Data Sumber: Clinical Lab Authority v1.0.4</p>
            </div>
          </aside>
        </div>
        
        {/* Article Card Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-1 bg-outline-variant/20 mb-1">
          <article className="bg-surface-container-lowest p-8 flex flex-col group transition-all duration-300 hover:bg-surface-container-low">
            <div className="mb-6 overflow-hidden aspect-[4/3] bg-surface-container">
              <img className="w-full h-full object-cover grayscale opacity-80 group-hover:scale-105 transition-transform duration-500" alt="Dopamine" src="https://images.unsplash.com/photo-1559757175-5700dde675bc?auto=format&fit=crop&q=80&w=800" />
            </div>
            <span className="font-label font-bold text-[10px] uppercase tracking-widest text-on-surface-variant mb-2">{t.knowledge.psychology}</span>
            <h3 className="font-headline font-bold text-xl tracking-tight uppercase mb-4 leading-tight text-on-surface">{t.knowledge.dopamineTitle}</h3>
            <div className="flex justify-between items-center pt-6 border-t border-outline-variant/10 mt-auto">
              <span className="material-symbols-outlined group-hover:translate-x-1 transition-transform">arrow_forward</span>
            </div>
          </article>
          
          <article className="bg-surface-container-lowest p-8 flex flex-col group transition-all duration-300 hover:bg-surface-container-low">
            <div className="mb-6 overflow-hidden aspect-[4/3] bg-surface-container">
              <img className="w-full h-full object-cover grayscale opacity-80 group-hover:scale-105 transition-transform duration-500" alt="RNG" src="https://images.unsplash.com/photo-1509228468518-180dd4864904?auto=format&fit=crop&q=80&w=800" />
            </div>
            <span className="font-label font-bold text-[10px] uppercase tracking-widest text-on-surface-variant mb-2">{t.knowledge.math}</span>
            <h3 className="font-headline font-bold text-xl tracking-tight uppercase mb-4 leading-tight text-on-surface">{t.knowledge.mythTitle}</h3>
            <div className="flex justify-between items-center pt-6 border-t border-outline-variant/10 mt-auto">
              <span className="material-symbols-outlined group-hover:translate-x-1 transition-transform">arrow_forward</span>
            </div>
          </article>
          
          <article className="bg-surface-container-lowest p-8 flex flex-col group transition-all duration-300 hover:bg-surface-container-low">
            <div className="mb-6 overflow-hidden aspect-[4/3] bg-surface-container">
              <img className="w-full h-full object-cover grayscale opacity-80 group-hover:scale-105 transition-transform duration-500" alt="Financial" src="https://images.unsplash.com/photo-1554224155-6726b3ff858f?auto=format&fit=crop&q=80&w=800" />
            </div>
            <span className="font-label font-bold text-[10px] uppercase tracking-widest text-on-surface-variant mb-2">{t.knowledge.finance}</span>
            <h3 className="font-headline font-bold text-xl tracking-tight uppercase mb-4 leading-tight text-on-surface">{t.knowledge.sunkCostTitle}</h3>
            <div className="flex justify-between items-center pt-6 border-t border-outline-variant/10 mt-auto">
              <span className="material-symbols-outlined group-hover:translate-x-1 transition-transform">arrow_forward</span>
            </div>
          </article>
        </div>

        {/* Dynamic Articles Grid */}
        {dynamicArticles.length > 0 && (
          <div className="mt-1">
            <div className="bg-surface-container-highest p-4 mb-1">
              <span className="font-label font-black text-[10px] uppercase tracking-[0.3em] text-primary">Latest From Clinical Authority</span>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-1 bg-outline-variant/20">
              {dynamicArticles.map(article => (
                <Link 
                  key={article.id} 
                  to={`/article/${article.id}`}
                  className="bg-surface-container-lowest p-8 flex flex-col group transition-all duration-300 hover:bg-surface-container-low"
                >
                  <span className="font-label font-bold text-[10px] uppercase tracking-widest text-primary mb-2">{article.category || 'Clinical Insight'}</span>
                  <h3 className="font-headline font-bold text-xl tracking-tight uppercase mb-4 leading-tight text-on-surface">{article.title}</h3>
                  <p className="font-body text-sm text-on-surface-variant line-clamp-3 mb-6">
                    {article.content}
                  </p>
                  <div className="flex justify-between items-center pt-6 border-t border-outline-variant/10 mt-auto">
                    <span className="font-label font-black text-[10px] uppercase tracking-widest group-hover:text-primary transition-colors">Read Full Insight</span>
                    <span className="material-symbols-outlined group-hover:translate-x-1 transition-transform">arrow_forward</span>
                  </div>
                </Link>
              ))}
            </div>
          </div>
        )}
      </main>
      <Footer />
    </div>
  );
}
