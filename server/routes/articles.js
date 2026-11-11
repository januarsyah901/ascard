import express from 'express';
import { prisma } from '../db.js';
import auth from '../middleware/auth.js';

const router = express.Router();

const pickArticle = (body = {}) => {
  const { title, content, category, published } = body;
  if (typeof title !== 'string' || !title.trim() || title.length > 200) return null;
  if (typeof content !== 'string' || !content.trim() || content.length > 100_000) return null;
  return {
    title: title.trim(),
    content,
    category: typeof category === 'string' && category.trim() ? category.trim().slice(0, 50) : null,
    published: published === true,
  };
};

const handlePrismaError = (error, res, next) => {
  if (error.code === 'P2025') return res.status(404).json({ message: 'Article not found' });
  if (error.code === 'P2002') return res.status(409).json({ message: 'An article with this title already exists' });
  next(error);
};

// All articles including drafts (Protected) — must be declared before "/:id"
router.get('/admin/all', auth, async (req, res, next) => {
  try {
    res.json(await prisma.article.findMany({ orderBy: { createdAt: 'desc' } }));
  } catch (error) {
    next(error);
  }
});

// Published articles (Public)
router.get('/', async (req, res, next) => {
  try {
    res.json(await prisma.article.findMany({
      where: { published: true },
      orderBy: { createdAt: 'desc' },
    }));
  } catch (error) {
    next(error);
  }
});

// Single published article (Public)
router.get('/:id', async (req, res, next) => {
  try {
    const article = await prisma.article.findFirst({
      where: { id: req.params.id, published: true },
    });
    if (!article) return res.status(404).json({ message: 'Article not found' });
    res.json(article);
  } catch (error) {
    next(error);
  }
});

// Create article (Protected)
router.post('/', auth, async (req, res, next) => {
  const data = pickArticle(req.body);
  if (!data) return res.status(400).json({ message: 'Title and content are required' });
  try {
    res.status(201).json(await prisma.article.create({ data }));
  } catch (error) {
    handlePrismaError(error, res, next);
  }
});

// Update article (Protected)
router.put('/:id', auth, async (req, res, next) => {
  const data = pickArticle(req.body);
  if (!data) return res.status(400).json({ message: 'Title and content are required' });
  try {
    res.json(await prisma.article.update({ where: { id: req.params.id }, data }));
  } catch (error) {
    handlePrismaError(error, res, next);
  }
});

// Delete article (Protected)
router.delete('/:id', auth, async (req, res, next) => {
  try {
    await prisma.article.delete({ where: { id: req.params.id } });
    res.json({ message: 'Article deleted successfully' });
  } catch (error) {
    handlePrismaError(error, res, next);
  }
});

export default router;
