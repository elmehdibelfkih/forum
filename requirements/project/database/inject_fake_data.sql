-- ====================================================================
-- seed_data.sql
-- Data‑only seeding for your existing forum.db schema
-- ====================================================================

PRAGMA foreign_keys = OFF;

-- 1) Delete all data (child → parent order)
DELETE FROM likes_dislikes;
DELETE FROM comments;
DELETE FROM post_categories;
DELETE FROM sessions;
DELETE FROM posts;
DELETE FROM categories_count;
DELETE FROM post_metadata;
DELETE FROM categories;
DELETE FROM users;

-- 2) Reset AUTOINCREMENT counters
DELETE FROM sqlite_sequence
 WHERE name IN (
   'users','categories','posts',
   'sessions','comments','likes_dislikes',
   'post_categories','post_metadata','categories_count'
 );

PRAGMA foreign_keys = ON;

-- 3) Seed users (IDs 1–3)
INSERT INTO users (username, email, password_hash) VALUES
  ('testuser','test@example.com','hash1'),
  ('alice',   'alice@example.com',   'hash2'),
  ('bob',     'bob@example.com',     'hash3');

-- 4) Seed 12 categories
INSERT INTO categories (name) VALUES
  ('Software Engineering'),
  ('Artificial Intelligence and Machine Learning'),
  ('Data Science and Big Data'),
  ('Cybersecurity'),
  ('Networking and Telecommunications'),
  ('Cloud Computing and Virtualization'),
  ('DevOps and SRE'),
  ('Database Systems'),
  ('Systems Programming'),
  ('Reverse Engineering'),
  ('Mobile and Embedded Development'),
  ('IoT (Internet of Things)');

-- 5) Seed 12 posts (user_id=1), titles = category names, long content
INSERT INTO posts (user_id, title, content) VALUES
  (1,'Software Engineering',
    'Software engineering is the systematic application of engineering principles to the development, operation, and maintenance of software systems. It encompasses the full lifecycle of software creation, from gathering requirements and designing systems to coding, testing, deploying, and maintaining applications. Software engineers work across a spectrum of programming paradigms and technologies to deliver solutions that are scalable, reliable, and maintainable.
At its core, software engineering is about problem-solving. Engineers must understand not only the technical requirements but also the business context and user needs. This discipline integrates aspects of project management, user experience design, system architecture, quality assurance, and documentation. Key methodologies include Agile, DevOps, and Waterfall, with each offering a structured approach to delivering software.
Software engineering is divided into various subfields, including front-end and back-end development, full-stack development, embedded systems, and cloud-based software. As software systems become more complex, the emphasis on architecture, modularity, and formal verification grows. Modern trends include microservices, containerization (e.g., Docker), and continuous integration/continuous deployment (CI/CD). Software engineers often rely on version control systems like Git and build tools such as Maven or Gradle.
Due to the increasing pervasiveness of software in all sectors, software engineering is one of the most vital and rapidly evolving fields in technology. Its principles are foundational to all other branches of computing.'
  ),


  (1,'Artificial Intelligence and Machine Learning',
    'AI and ML focus on building systems that learn from data. From training neural networks for image recognition to \
reinforcement learning agents, the tech spans data preprocessing, model selection, hyperparameter tuning, and MLOps pipelines.'
  ),


  (1,'Data Science and Big Data',
    'Data Science and Big Data involve collecting, processing, and analyzing massive datasets. Tools like Hadoop, Spark, \
and lakehouse architectures enable scalable ETL, analytics, and real‑time insights, while governance and privacy remain critical.'
  ),


  (1,'Cybersecurity',
    'Cybersecurity protects systems and data from threats. It covers identity management, encryption, network and endpoint \
security, threat intelligence, and incident response, guided by frameworks like Zero Trust and NIST.'
  ),


  (1,'Networking and Telecommunications',
    'This field underpins digital communication—routing algorithms, TCP/IP, wireless (Wi‑Fi, 5G), optical networks, and new \
paradigms like SCION and Named Data Networking for better security and performance.'
  ),


  (1,'Cloud Computing and Virtualization',
    'Cloud and virtualization provide on‑demand compute, storage, and networking. VMs, containers, and serverless functions \
are orchestrated by platforms like Kubernetes to deliver resilient, scalable applications.'
  ),


  (1,'DevOps and SRE',
    'DevOps and SRE bridge dev and ops. DevOps uses infrastructure as code and CI/CD for rapid delivery, while SRE defines \
Service Level Objectives, error budgets, and blameless postmortems to maintain reliability.'
  ),


  (1,'Database Systems',
    'Database Systems store and retrieve data efficiently. Relational databases offer ACID and SQL, NoSQL handles unstructured \
workloads, and modern variants like NewSQL and vector databases serve AI‑driven use cases.'
  ),


  (1,'Systems Programming',
    'Systems Programming writes low‑level code interacting with hardware and kernels. Languages such as C, Rust, and Zig provide \
different safety, performance, and abstraction trade‑offs for firmware, drivers, and OS components.'
  ),


  (1,'Reverse Engineering',
    'Reverse Engineering analyzes binaries or hardware to reveal internal logic. Tools like Ghidra, IDA Pro, and Radare2 help \
unpack obfuscated code, perform dynamic tracing, and discover vulnerabilities or compatibility details.'
  ),

  
  (1,'Mobile and Embedded Development',
    'Mobile and embedded targets range from smartphones to microcontrollers. Frameworks (Android SDK, iOS Swift) and RTOS \
(FreeRTOS, Zephyr) must balance UI/UX, resource constraints, power management, and OTA update pipelines.'
  ),


  (1,'IoT (Internet of Things)',
    'IoT connects sensors and actuators at the edge. Lightweight protocols (MQTT, CoAP), secure provisioning (TPM/HSM), and \
edge computing combine to support scalable, data‑driven automation across industries.'
  );

-- 6) Link each post to its matching category
INSERT INTO post_categories (post_id, category_id)
SELECT p.id, c.id
  FROM posts p
  JOIN categories c ON p.title = c.name;

-- 7) Seed sessions
INSERT INTO sessions (user_id, session_token, expires_at) VALUES
  (1,'sess1',DATETIME('now','+1 day')),
  (2,'sess2',DATETIME('now','+2 days')),
  (3,'sess3',DATETIME('now','+3 days'));

-- 8) Add one comment per post (rotate users 2 & 3)
INSERT INTO comments (user_id, post_id, comment) VALUES
  (2,1,'Insightful overview of software engineering!'),
  (3,2,'Great breakdown of AI/ML practices.'),
  (2,3,'Big Data pipelines are so critical these days.'),
  (3,4,'Security is foundational, thanks for detailing it.'),
  (2,5,'Networking innovations keep us all connected.'),
  (3,6,'Kubernetes is a game‑changer for ops.'),
  (2,7,'Error budgets really do improve reliability.'),
  (3,8,'The right DB choice is everything.'),
  (2,9,'Systems code safety is a fascinating topic.'),
  (3,10,'Reversing binaries always feels like detective work.'),
  (2,11,'Embedded constraints make development fun.'),
  (3,12,'IoT security at scale remains a big challenge.');

-- 9) Add likes/dislikes
INSERT INTO likes_dislikes (user_id, post_id, is_like, is_dislike) VALUES
  -- Alice (user 2) likes posts 1–6, dislikes 7
  (2,1,1,0),(2,2,1,0),(2,3,1,0),(2,4,1,0),(2,5,1,0),(2,6,1,0),(2,7,0,1),
  -- Bob (user 3) likes posts 7–12, dislikes 1
  (3,7,1,0),(3,8,1,0),(3,9,1,0),(3,10,1,0),(3,11,1,0),(3,12,1,0),(3,1,0,1);

-- 10) Populate post_metadata (total posts = 12)
DELETE FROM post_metadata;
INSERT INTO post_metadata (id, post_count) VALUES (1, 12);

-- 11) Populate categories_count (each = 1)
DELETE FROM categories_count;
INSERT INTO categories_count (category_id, post_count) VALUES
  (1,1),(2,1),(3,1),(4,1),(5,1),(6,1),
  (7,1),(8,1),(9,1),(10,1),(11,1),(12,1);
