import express from 'express';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import { prisma } from '../db.js';

const router = express.Router();
const JWT_SECRET = process.env.JWT_SECRET || 'fallback-secret';

// Admin Login
router.post('/login', async (req, res) => {
  const { username, password } = req.body;

  try {
    const admin = await prisma.admin.findUnique({
      where: { username }
    });

    if (!admin) {
      return res.status(401).json({ message: 'Invalid credentials' });
    }

    const isMatch = await bcrypt.compare(password, admin.password);
    if (!isMatch) {
      return res.status(401).json({ message: 'Invalid credentials' });
    }

    const token = jwt.sign({ id: admin.id }, JWT_SECRET, { expiresIn: '1d' });
    res.json({ token, admin: { id: admin.id, username: admin.username } });
  } catch (error) {
    res.status(500).json({ message: 'Server error', error: error.message });
  }
});

// Seed Initial Admin (for testing)
router.post('/seed', async (req, res) => {
  const { username, password } = req.body;
  const hashedPassword = await bcrypt.hash(password, 10);
  
  try {
    const admin = await prisma.admin.create({
      data: { username, password: hashedPassword }
    });
    res.json({ message: 'Admin seeded successfully', admin: { id: admin.id, username: admin.username } });
  } catch (error) {
    res.status(400).json({ message: 'Seed failed', error: error.message });
  }
});

export default router;
