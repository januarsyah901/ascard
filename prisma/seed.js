import 'dotenv/config';
import { prisma } from '../server/db.js';
import bcrypt from 'bcryptjs';

async function main() {
  console.log('Starting seeding...');

  // Create Admin
  const hashedPassword = await bcrypt.hash('password123', 10);
  const admin = await prisma.admin.upsert({
    where: { username: 'admin' },
    update: {},
    create: {
      username: 'admin',
      password: hashedPassword,
    },
  });
  console.log('Admin created:', admin.username);

  // Create Sample Articles
  const articles = [
    {
      title: 'Anatomi Jebakan Dopamin',
      content: 'Dopamin bukan hormon kebahagiaan, melainkan hormon antisipasi. Di balik setiap putaran mesin slot, ada lonjakan dopamin yang memicu rasa penasaran berlebih. Secara klinis, sistem reward otak dibajak oleh ketidakpastian yang dirancang secara matematis.',
      category: 'Psikologi',
      published: true,
    },
    {
      title: 'Mitos Keberuntungan dalam RNG',
      content: 'Random Number Generator (RNG) tidak mengenal pola. Setiap putaran bersifat independen. Mempercayai "pola" adalah awal dari kekalahan finansial. Secara statistik, probabilitas tidak memiliki memori.',
      category: 'Matematika',
      published: true,
    },
    {
      title: 'Biaya Tersembunyi: Sunk Cost Fallacy',
      content: 'Semakin banyak uang yang Anda keluarkan, semakin sulit Anda berhenti. Ini adalah jebakan psikologis yang sering digunakan dalam desain kasino. Otak kita membenci kehilangan, sehingga kita terus "berinvestasi" pada kekalahan.',
      category: 'Finansial',
      published: true,
    }
  ];

  for (const article of articles) {
    await prisma.article.upsert({
      where: { title: article.title },
      update: {},
      create: article,
    });
  }

  console.log('Seeding finished successfully.');
}

main()
  .catch((e) => {
    console.error(e);
    process.exit(1);
  })
  .finally(async () => {
    await prisma.$disconnect();
  });
