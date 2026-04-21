import React from 'react';
import Navigation from '../components/Navigation';
import Footer from '../components/Footer';
import { Link } from 'react-router-dom';
import { useLanguage } from '../context/LanguageContext';
import { useCurrency } from '../context/CurrencyContext';

export default function LandingPage() {
  const { t } = useLanguage();
  const { formatValue } = useCurrency();
  
  return (
    <>
      <Navigation />
      <main className="pt-16">
        {/* Hero Section */}
        <section className="relative min-h-[921px] flex items-center px-6 md:px-12 bg-surface">
          <div className="grid grid-cols-1 lg:grid-cols-12 w-full max-w-7xl mx-auto gap-12 py-20">
            <div className="lg:col-span-8 flex flex-col justify-center">
              <div className="inline-block px-3 py-1 bg-primary text-on-primary text-[10px] font-bold tracking-[0.2em] uppercase mb-6 self-start">
                {t.landing.researchLab}
              </div>
              <h1 className="text-6xl md:text-8xl lg:text-9xl font-black tracking-tighter text-primary leading-[0.85] uppercase mb-8">
                {t.landing.heroTitle}
              </h1>
              <p className="text-lg md:text-xl text-on-surface-variant max-w-2xl font-light leading-relaxed mb-10 border-l-4 border-primary pl-6">
                {t.landing.heroDesc}
              </p>
              <div className="flex flex-col sm:flex-row gap-4">
                <Link to="/simulation">
                  <button className="monolith-gradient px-10 py-5 text-on-primary font-bold tracking-tighter uppercase font-headline text-sm hover:opacity-90 transition-opacity w-full">
                    {t.landing.ctaSimulation}
                  </button>
                </Link>
                <Link to="/knowledge-base">
                  <button className="px-10 py-5 bg-transparent border-2 border-primary text-primary font-bold tracking-tighter uppercase font-headline text-sm hover:bg-surface-container transition-colors w-full">
                    {t.landing.ctaMethodology}
                  </button>
                </Link>
              </div>
            </div>
            <div className="lg:col-span-4 relative flex items-center justify-center">
              <div className="w-full aspect-square bg-surface-container-highest relative overflow-hidden">
                <img className="w-full h-full object-cover grayscale mix-blend-multiply opacity-80" alt="abstract high-contrast digital grid" src="https://lh3.googleusercontent.com/aida-public/AB6AXuARIdrAG5_Z3f9WjRKuFOMi-W5CopaPGAZhrtvCFO3bUa8-rtyFKNzCqiGptX_thZjk69U1X4seioXsGke0wLvw4FaeQaWrsP_m-56qdaXC4e02xkKNx6y9O6HHXtymjJ9S-PkjYVGFRMZ3DaaIX1Wrc3e5VsOY-mDiMRDW1KpYcC0n231KeR-01zek-V7HpM-REu8ifv8xDEm9gIJ86-EAFZuerftxjs9sjBipyQ6UgP1Xpz_5SVcPxIC98-Bkzz52A4gB2ccLofkH" />
                <div className="absolute inset-0 bg-gradient-to-t from-surface to-transparent"></div>
              </div>
              <div className="absolute -bottom-6 -left-6 bg-primary p-6 text-on-primary w-48 hidden md:block">
                <p className="text-[10px] tracking-widest uppercase mb-2 opacity-60">{t.landing.systemStatus}</p>
                <p className="text-xs font-bold leading-tight uppercase font-headline">{t.landing.activeEncrypted}</p>
              </div>
            </div>
          </div>
        </section>

        {/* Stats Section (Empirical Monolith Bento Grid) */}
        <section className="py-24 px-6 md:px-12 bg-surface-container-low">
          <div className="max-w-7xl mx-auto">
            <div className="mb-16">
              <h2 className="text-sm font-bold tracking-[0.3em] uppercase text-outline mb-2">{t.landing.empiricalData}</h2>
              <div className="h-1 w-20 bg-primary"></div>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-1">
              {/* Main Stat Card */}
              <div className="md:col-span-2 bg-surface p-12 flex flex-col justify-between min-h-[400px] border border-outline-variant/10">
                <div>
                  <span className="material-symbols-outlined text-4xl mb-6 block text-primary">analytics</span>
                  <h3 className="text-3xl font-bold tracking-tighter uppercase mb-4">{t.landing.probTitle}</h3>
                </div>
                <div className="flex flex-col gap-2">
                  <span className="text-7xl md:text-9xl font-black text-primary tracking-tighter leading-none">{t.landing.probValue}</span>
                  <p className="text-xl md:text-2xl font-medium text-on-surface uppercase tracking-tight max-w-md">
                    {t.landing.probDesc}
                  </p>
                </div>
              </div>
              {/* Side Stat Column */}
              <div className="flex flex-col gap-1">
                <div className="bg-primary p-8 flex flex-col justify-center flex-1">
                  <p className="text-on-primary-fixed-variant text-[10px] font-bold tracking-widest uppercase mb-2">{t.landing.houseEdge}</p>
                  <p className="text-4xl font-black text-on-primary tracking-tighter mb-4">97.4%</p>
                  <p className="text-xs text-on-primary-fixed-variant leading-relaxed opacity-70">
                    {t.landing.houseEdgeDesc}
                  </p>
                </div>
                <div className="bg-surface-container-highest p-8 flex flex-col justify-center flex-1">
                  <span className="material-symbols-outlined text-primary mb-4" style={{ fontVariationSettings: "'FILL' 1" }}>trending_down</span>
                  <p className="text-2xl font-bold text-primary tracking-tighter uppercase">{t.landing.cognitiveTitle}</p>
                  <p className="text-xs text-on-surface-variant leading-relaxed mt-2 uppercase font-medium">
                    {t.landing.cognitiveDesc}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Reality Check Section */}
        <section className="py-32 px-6 md:px-12 bg-surface">
          <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center gap-16">
            <div className="w-full md:w-1/2">
              <img className="w-full grayscale border-l-[12px] border-primary" alt="microscopic detail" src="https://lh3.googleusercontent.com/aida-public/AB6AXuCCdPsqVrsaGPSWnqjapXd_GCFL9aXSSOmD9R-U7SAqxQiXldlrtzifnAVRfIVb1Jbro03-Dns1I1ACYvFVxSh7KcB74u5IWWWaV87gxmmUi9VZmcV4T1bdv7lmWdGSGI8bkqQKzOLxpVP1udg98rFJQ-vdZk3IjothJC43-ANepxOoLxn3UtbzF-eow4mO-SiH2KH1GQjZ77lw2TKEzA8IYXDEDUFsithnAVJJhLzV6Z4hxnAXCoHKIQL7F_FQh5HWe01leS58RymE" />
            </div>
            <div className="w-full md:w-1/2">
              <h2 className="text-5xl font-black tracking-tighter uppercase mb-8 leading-none">{t.landing.whyFooled}</h2>
              <div className="space-y-8">
                <div className="group">
                  <div className="flex items-start gap-4 mb-2">
                    <span className="text-sm font-black text-primary bg-surface-container px-2 py-1">01</span>
                    <h4 className="text-lg font-bold uppercase tracking-tight">{t.landing.point1Title}</h4>
                  </div>
                  <p className="text-on-surface-variant text-sm leading-relaxed pl-10">
                    {t.landing.point1Desc}
                  </p>
                </div>
                <div className="group">
                  <div className="flex items-start gap-4 mb-2">
                    <span className="text-sm font-black text-primary bg-surface-container px-2 py-1">02</span>
                    <h4 className="text-lg font-bold uppercase tracking-tight">{t.landing.point2Title}</h4>
                  </div>
                  <p className="text-on-surface-variant text-sm leading-relaxed pl-10">
                    {t.landing.point2Desc}
                  </p>
                </div>
                <div className="group">
                  <div className="flex items-start gap-4 mb-2">
                    <span className="text-sm font-black text-primary bg-surface-container px-2 py-1">03</span>
                    <h4 className="text-lg font-bold uppercase tracking-tight">{t.landing.point3Title}</h4>
                  </div>
                  <p className="text-on-surface-variant text-sm leading-relaxed pl-10">
                    {t.landing.point3Desc.replace('{win}', formatValue(0.50)).replace('{bet}', formatValue(1.00))}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* Final CTA Monolith */}
        <section className="py-24 px-6 md:px-12 bg-primary">
          <div className="max-w-5xl mx-auto text-center">
            <h2 className="text-4xl md:text-6xl font-black tracking-tighter text-on-primary uppercase mb-12">
              {t.landing.finalTitle}
            </h2>
            <div className="flex flex-col md:flex-row justify-center gap-6">
              <Link to="/simulation">
                <button className="bg-surface text-primary px-12 py-6 font-black tracking-widest uppercase text-base hover:bg-surface-dim transition-colors w-full">
                  {t.landing.finalDashboard}
                </button>
              </Link>
              <Link to="/intervention">
                <button className="border-2 border-on-primary text-on-primary px-12 py-6 font-black tracking-widest uppercase text-base hover:bg-on-primary hover:text-primary transition-colors w-full">
                  {t.landing.finalConsult}
                </button>
              </Link>
            </div>
            <p className="mt-12 text-[10px] text-on-primary/40 tracking-[0.4em] uppercase font-bold">
              DATA DRIVEN • CLINICALLY VALIDATED • NON-PROFIT MISSION
            </p>
          </div>
        </section>
      </main>
      <Footer />
    </>
  );
}
