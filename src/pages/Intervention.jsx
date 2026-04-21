import React from 'react';
import Navigation from '../components/Navigation';
import Footer from '../components/Footer';
import { useLanguage } from '../context/LanguageContext';

export default function Intervention() {
  const { t } = useLanguage();

  return (
    <div className="min-h-screen flex flex-col bg-surface">
      <Navigation />
      <main className="flex-grow pt-24 pb-20 px-6 md:px-12 max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-12">
        {/* LEFT COLUMN: HEADER & FAQ ACCORDIONS */}
        <section className="lg:col-span-7 space-y-12">
          <header>
            <span className="label-sm text-[10px] uppercase tracking-[0.2em] text-secondary font-headline font-bold block mb-2">{t.intervention.empiricalDeconstruction}</span>
            <h1 className="text-6xl md:text-8xl font-black tracking-tighter leading-none uppercase mb-6">{t.intervention.mythLab.split(' ')[0]} <br /> {t.intervention.mythLab.split(' ')[1]}</h1>
            <p className="text-secondary max-w-lg leading-relaxed text-lg">{t.intervention.mythDesc}</p>
          </header>
          <div className="space-y-4">
            {/* Accordion Item 1 */}
            <div className="bg-surface-container-low border-b border-outline-variant/30">
              <details className="group" open>
                <summary className="flex justify-between items-center p-6 cursor-pointer list-none">
                  <span className="font-headline font-extrabold text-xl tracking-tight uppercase">{t.intervention.myth1Title}</span>
                  <span className="material-symbols-outlined transition-transform group-open:rotate-180">expand_more</span>
                </summary>
                <div className="px-6 pb-8 text-on-surface">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div className="space-y-4">
                      <h4 className="font-bold text-sm uppercase tracking-widest text-secondary">The Myth</h4>
                      <p className="text-sm">{t.intervention.myth1Myth}</p>
                    </div>
                    <div className="space-y-4 border-l border-outline-variant/50 pl-8">
                      <h4 className="font-bold text-sm uppercase tracking-widest text-primary">Mathematical Fact</h4>
                      <p className="text-sm italic">{t.intervention.myth1Fact}</p>
                    </div>
                  </div>
                </div>
              </details>
            </div>
            
            {/* Accordion Item 2 */}
            <div className="bg-surface-container-low border-b border-outline-variant/30">
              <details className="group">
                <summary className="flex justify-between items-center p-6 cursor-pointer list-none">
                  <span className="font-headline font-extrabold text-xl tracking-tight uppercase">{t.intervention.myth2Title}</span>
                  <span className="material-symbols-outlined transition-transform group-open:rotate-180">expand_more</span>
                </summary>
                <div className="px-6 pb-8 text-on-surface">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div className="space-y-4">
                      <h4 className="font-bold text-sm uppercase tracking-widest text-secondary">The Myth</h4>
                      <p className="text-sm">{t.intervention.myth2Myth}</p>
                    </div>
                    <div className="space-y-4 border-l border-outline-variant/50 pl-8">
                      <h4 className="font-bold text-sm uppercase tracking-widest text-primary">Mathematical Fact</h4>
                      <p className="text-sm italic">{t.intervention.myth2Fact}</p>
                    </div>
                  </div>
                </div>
              </details>
            </div>
            
            {/* Accordion Item 3 */}
            <div className="bg-surface-container-low border-b border-outline-variant/30">
              <details className="group">
                <summary className="flex justify-between items-center p-6 cursor-pointer list-none">
                  <span className="font-headline font-extrabold text-xl tracking-tight uppercase">{t.intervention.myth3Title}</span>
                  <span className="material-symbols-outlined transition-transform group-open:rotate-180">expand_more</span>
                </summary>
                <div className="px-6 pb-8 text-on-surface">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <div className="space-y-4">
                      <h4 className="font-bold text-sm uppercase tracking-widest text-secondary">The Myth</h4>
                      <p className="text-sm">{t.intervention.myth3Myth}</p>
                    </div>
                    <div className="space-y-4 border-l border-outline-variant/50 pl-8">
                      <h4 className="font-bold text-sm uppercase tracking-widest text-primary">Mathematical Fact</h4>
                      <p className="text-sm italic">{t.intervention.myth3Fact}</p>
                    </div>
                  </div>
                </div>
              </details>
            </div>
          </div>
        </section>
        
        {/* RIGHT COLUMN: AI CHAT INTERVENTION */}
        <aside className="lg:col-span-5 relative">
          <div className="sticky top-24 bg-surface-container-lowest p-8 border border-outline-variant/20 shadow-2xl">
            <div className="flex items-center justify-between mb-8">
              <div>
                <h3 className="font-headline font-black text-2xl uppercase tracking-tighter">{t.intervention.logicModule}</h3>
                <p className="text-[10px] font-mono text-secondary">{t.intervention.cbtVersion}</p>
              </div>
              <div className="w-3 h-3 bg-primary animate-pulse"></div>
            </div>
            <div className="h-[400px] overflow-y-auto mb-6 flex flex-col gap-4 no-scrollbar">
              <div className="bg-surface-container p-4 max-w-[85%] self-start">
                <p className="text-xs leading-relaxed">{t.intervention.aiWelcome}</p>
              </div>
              <div className="bg-primary p-4 max-w-[85%] self-end">
                <p className="text-xs text-on-primary leading-relaxed">{t.intervention.userExample}</p>
              </div>
              <div className="bg-surface-container p-4 max-w-[85%] self-start border-l-4 border-primary">
                <p className="text-xs leading-relaxed font-bold uppercase mb-2">{t.intervention.aiCorrection}</p>
                <p className="text-xs leading-relaxed italic">{t.intervention.aiResponse}</p>
              </div>
            </div>
            <div className="relative border-t-2 border-primary pt-4">
              <input className="w-full bg-transparent border-none focus:ring-0 placeholder:text-neutral-300 font-headline uppercase text-sm px-0" placeholder={t.intervention.queryPlaceholder} type="text" />
              <button className="absolute right-0 top-4">
                <span className="material-symbols-outlined text-primary">send</span>
              </button>
            </div>
          </div>
        </aside>
      </main>
      <Footer />
    </div>
  );
}
