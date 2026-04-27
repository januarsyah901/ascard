import express from 'express';
import { prisma } from '../db.js';
import auth from '../middleware/auth.js';

const router = express.Router();

// Get all articles (Public)
router.get('/', async (req, res) => {
  try {
    const articles = await prisma.article.findMany({
      orderBy: { createdAt: 'desc' }
    });
    res.json(articles);
  } catch (error) {
    res.status(500).json({ message: 'Error fetching articles', error: error.message });
  }
});

// Get single article (Public)
router.get('/:id', async (req, res) => {
  const { id } = req.params;
  try {
    const article = await prisma.article.findUnique({
      where: { id }
    });
    if (!article) return res.status(404).json({ message: 'Article not found' });
    res.json(article);
  } catch (error) {
    res.status(500).json({ message: 'Error fetching article', error: error.message });
  }
});

// Create article (Protected)
router.post('/', auth, async (req, res) => {
  const { title, content, category, published } = req.body;
  try {
    const article = await prisma.article.create({
      data: { title, content, category, published }
    });
    res.json(article);
  } catch (error) {
    res.status(400).json({ message: 'Error creating article', error: error.message });
  }
});

// Update article (Protected)
router.put('/:id', auth, async (req, res) => {
  const { id } = req.params;
  const { title, content, category, published } = req.body;
  try {
    const article = await prisma.article.update({
      where: { id },
      data: { title, content, category, published }
    });
    res.json(article);
  } catch (error) {
    res.status(400).json({ message: 'Error updating article', error: error.message });
  }
});

// Delete article (Protected)
router.delete('/:id', auth, async (req, res) => {
  const { id } = req.params;
  try {
    await prisma.article.delete({ where: { id } });
    res.json({ message: 'Article deleted successfully' });
  } catch (error) {
    res.status(400).json({ message: 'Error deleting article', error: error.message });
  }
});

export default router;
