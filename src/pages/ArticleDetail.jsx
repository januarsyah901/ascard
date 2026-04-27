import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import Navigation from '../components/Navigation';
import Footer from '../components/Footer';
import { ArrowLeft, Clock, Tag, Share2 } from 'lucide-react';

const ArticleDetail = () => {
  const { id } = useParams();
  const [article, setArticle] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:5001/api/articles/${id}`)
      .then(res => res.json())
      .then(data => {
        setArticle(data);
        setLoading(false);
      })
      .catch(err => {
        console.error('Error fetching article:', err);
        setLoading(false);
      });
  }, [id]);

  if (loading) {
    return (
      <div className="min-h-screen bg-surface flex items-center justify-center">
        <div className="animate-pulse font-label font-bold text-[10px] uppercase tracking-[0.4em]">Analyzing Data...</div>
      </div>
    );
  }

  if (!article) {
    return (
      <div className="min-h-screen bg-surface flex flex-col items-center justify-center p-6 text-center">
        <h1 className="font-headline font-black text-4xl uppercase mb-4">Article Not Found</h1>
        <Link to="/knowledge-base" className="font-label font-bold text-xs uppercase tracking-widest border-b-2 border-primary pb-1">Return to Repository</Link>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex flex-col bg-surface text-on-surface font-body">
      <Navigation />
      
      <main className="flex-grow pt-32 pb-24 px-6 md:px-12 max-w-4xl mx-auto w-full">
        <Link 
          to="/knowledge-base" 
          className="inline-flex items-center gap-2 text-outline hover:text-primary transition-colors mb-12 group"
        >
          <ArrowLeft className="w-4 h-4 group-hover:-translate-x-1 transition-transform" />
          <span className="font-label font-bold text-[10px] uppercase tracking-widest">Back to Insights</span>
        </Link>

        <article className="animate-in fade-in slide-in-from-bottom-4 duration-700">
          <header className="mb-12">
            <div className="flex items-center gap-4 mb-6">
              <span className="bg-primary text-on-primary-fixed px-3 py-1 font-label text-[10px] font-bold tracking-widest uppercase">
                {article.category || 'Clinical Insight'}
              </span>
              <div className="flex items-center gap-2 text-outline font-label text-[10px] uppercase tracking-widest">
                <Clock className="w-3 h-3" />
                <span>{new Date(article.createdAt).toLocaleDateString()}</span>
              </div>
            </div>
            
            <h1 className="font-headline font-black text-5xl md:text-7xl tracking-tighter uppercase leading-[0.9] mb-8">
              {article.title}
            </h1>

            <div className="h-1 w-20 bg-primary"></div>
          </header>

          <div className="prose prose-neutral max-w-none">
            <div className="font-body text-xl leading-relaxed text-on-surface-variant whitespace-pre-wrap">
              {article.content}
            </div>
          </div>

          <footer className="mt-20 pt-12 border-t border-outline-variant/30">
            <div className="flex flex-wrap items-center justify-between gap-6">
              <div className="flex items-center gap-3">
                <Tag className="w-4 h-4 text-outline" />
                <span className="font-label font-bold text-[10px] uppercase tracking-widest text-outline">
                  Keywords: {article.category}, Clinical, Behavioral, RNG
                </span>
              </div>
              <button className="flex items-center gap-2 bg-surface-container-high px-4 py-2 hover:bg-surface-container-highest transition-all group">
                <Share2 className="w-4 h-4 text-primary" />
                <span className="font-label font-bold text-[10px] uppercase tracking-widest">Share Analysis</span>
              </button>
            </div>
          </footer>
        </article>
      </main>

      <Footer />
    </div>
  );
};

export default ArticleDetail;
