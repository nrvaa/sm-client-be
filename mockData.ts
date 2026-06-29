export interface TimelineEvent {
  id: string;
  status: string; // e.g., "Disassembly", "Body & Paint", "Engine & Drivetrain", "Interior", "Final Detailing", "Road Test"
  date: string;
  description: string;
  image?: string;
  completed: boolean;
}

export interface ProgressPhoto {
  id: string;
  url: string;
  title: string;
  date: string;
  stage: string;
  type?: "photo" | "video";
}

export interface ReplacedPart {
  id: string;
  name: string;
  from: string;
  to: string;
}

export interface Vehicle {
  id: string;
  clientCode: string;
  clientName: string;
  brand: string;
  model: string;
  year: number;
  licensePlate: string;
  vin: string;
  status: string;
  completionPercentage: number;
  estimatedCompletion: string;
  bannerImage: string;
  restorationType: string; // e.g., "Restorasi Total" atau "Restorasi Parsial"
  timeline: TimelineEvent[];
  gallery: ProgressPhoto[];
  replacedParts?: ReplacedPart[];
}

export const initialVehicles: Vehicle[] = [
  {
    id: "v1",
    clientCode: "SM-BUDI",
    clientName: "Budi Santoso",
    brand: "Porsche",
    model: "911 Carrera RS 2.7",
    year: 1973,
    licensePlate: "B 911 RS",
    vin: "9113600123",
    status: "Paint Prep & Rust Repair",
    completionPercentage: 65,
    estimatedCompletion: "2026-10-15",
    bannerImage: "https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e?auto=format&fit=crop&q=80&w=1200",
    restorationType: "Restorasi Total",
    replacedParts: [
      { id: "rp1-1", name: "Jok Mobil", from: "Jok Kain Bawaan", to: "Jok Kulit Asli dengan Motif Houndstooth" },
      { id: "rp1-2", name: "Knalpot", from: "Knalpot Standar Karatan", to: "Knalpot Stainless Steel Dual Exhaust" },
      { id: "rp1-3", name: "Suspensi", from: "Shockbreaker Usang", to: "Suspensi Coilover Sport" }
    ],
    timeline: [
      {
        id: "t1-1",
        status: "Initial Inspection & Intake",
        date: "2026-02-10",
        description: "Kendaraan tiba di workshop Stanley Marthin. Melakukan inspeksi menyeluruh pada bodi, mesin, kelistrikan, dan inventarisasi komponen asli.",
        image: "https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t1-2",
        status: "Complete Disassembly",
        date: "2026-03-01",
        description: "Proses pembongkaran seluruh bodi mobil hingga menyisakan sasis kosong (bare shell). Semua komponen dikelompokkan dan dicatat kelayakannya.",
        image: "https://images.unsplash.com/photo-1486006920555-c77dce18193b?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t1-3",
        status: "Sandblasting & Metalwork",
        date: "2026-04-15",
        description: "Pembersihan karat menggunakan metode sandblasting. Pengelasan bagian bodi yang keropos dan restorasi panel bodi asli menggunakan plat baja baru.",
        image: "https://images.unsplash.com/photo-1534088568595-a066f410bcda?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t1-4",
        status: "Paint Prep & Epoxy Primer",
        date: "2026-06-10",
        description: "Pengaplikasian lapisan epoxy primer anti-karat pertama pada seluruh permukaan logam sasis. Pengisian dempul tipis untuk meratakan permukaan bodi mobil.",
        image: "https://images.unsplash.com/photo-1619642751034-765dfdf7c58e?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t1-5",
        status: "Painting & Clear Coat",
        date: "2026-07-20",
        description: "Proses pengecatan warna legendaris Grand Prix White dengan aksen decal Viper Green khas Carrera RS. Dilanjutkan dengan 3 lapis clear coat premium.",
        completed: false
      },
      {
        id: "t1-6",
        status: "Engine & Drivetrain Assembly",
        date: "2026-08-30",
        description: "Pemasangan kembali mesin flat-six 2.7L setelah restorasi total (rebuild piston, penggantian oli seal, restorasi pompa MFI).",
        completed: false
      },
      {
        id: "t1-7",
        status: "Interior Restoration & Fitment",
        date: "2026-09-20",
        description: "Pemasangan jok bucket klasik dengan bahan fabric bermotif houndstooth hitam-putih, setir orisinil, panel dashboard kulit baru, dan karpet velour hitam.",
        completed: false
      },
      {
        id: "t1-8",
        status: "Road Test & Final Detailing",
        date: "2026-10-10",
        description: "Pengujian jalan sejauh 100 km untuk kalibrasi suspensi, kelistrikan, dan mesin. Diakhiri dengan poles cat tingkat tinggi multi-stage detailing.",
        completed: false
      }
    ],
    gallery: [
      {
        id: "g1-1",
        url: "https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800",
        title: "Tiba di Stanley Marthin",
        date: "2026-02-10",
        stage: "Intake",
        type: "photo"
      },
      {
        id: "g1-2",
        url: "https://images.unsplash.com/photo-1486006920555-c77dce18193b?auto=format&fit=crop&q=80&w=800",
        title: "Pembongkaran Sasis & Mesin",
        date: "2026-03-05",
        stage: "Disassembly",
        type: "photo"
      },
      {
        id: "g1-3",
        url: "https://images.unsplash.com/photo-1534088568595-a066f410bcda?auto=format&fit=crop&q=80&w=800",
        title: "Pengelasan Panel Bodi Keropos",
        date: "2026-04-20",
        stage: "Metalwork",
        type: "photo"
      },
      {
        id: "g1-video",
        url: "https://www.stanleymarthin.com/wp-content/uploads/2018/04/StanleyMarthin2.mp4",
        title: "Proses Sandblasting Sasis",
        date: "2026-05-02",
        stage: "Metalwork",
        type: "video"
      },
      {
        id: "g1-4",
        url: "https://images.unsplash.com/photo-1619642751034-765dfdf7c58e?auto=format&fit=crop&q=80&w=800",
        title: "Epoxy Primer & Undercoat",
        date: "2026-06-15",
        stage: "Paint Prep",
        type: "photo"
      },
      {
        id: "g1-5",
        url: "https://images.unsplash.com/photo-1614162692292-7ac56d7f7f1e?auto=format&fit=crop&q=80&w=800",
        title: "Detailing Restorasi Kabin",
        date: "2026-06-25",
        stage: "Interior",
        type: "photo"
      }
    ]
  },
  {
    id: "v2",
    clientCode: "SM-BUDI",
    clientName: "Budi Santoso",
    brand: "Ford",
    model: "Mustang Fastback GT500",
    year: 1967,
    licensePlate: "DK 1967 GT",
    vin: "7T02S123456",
    status: "Suspension & Drivetrain Install",
    completionPercentage: 80,
    estimatedCompletion: "2026-08-25",
    bannerImage: "https://images.unsplash.com/photo-1615906655593-ad0386982a0f?auto=format&fit=crop&q=80&w=1200",
    restorationType: "Restorasi Parsial",
    replacedParts: [
      { id: "rp2-1", name: "Velg", from: "Velg Kaleng Standar", to: "Velg Racing Klasik 17 Inch" },
      { id: "rp2-2", name: "Rem", from: "Rem Tromol Belakang", to: "Rem Cakram Wilwood 4 Roda" }
    ],
    timeline: [
      {
        id: "t2-1",
        status: "Intake & Body Inspection",
        date: "2026-01-05",
        description: "Mobil Mustang Fastback tiba dalam kondisi karat berat di bagian lantai bawah dan bagasi. Semua komponen aksesoris diinventarisasi.",
        image: "https://images.unsplash.com/photo-1583121274602-3e2820c69888?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t2-2",
        status: "Chassis Sandblasting & Repair",
        date: "2026-02-15",
        description: "Restorasi penuh struktur lantai (floor pans) dan firewall menggunakan lembaran besi orisinil Mustang Dynacorn.",
        image: "https://images.unsplash.com/photo-1605558202076-168223c9f02b?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t2-3",
        status: "Pengecatan Base & Clear",
        date: "2026-04-10",
        description: "Warna Eleanor Gray Metallic dengan garis ganda hitam (stripes) selesai dikerjakan di oven cat profesional.",
        image: "https://images.unsplash.com/photo-1615906655593-ad0386982a0f?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t2-4",
        status: "Suspension & Drivetrain Install",
        date: "2026-06-02",
        description: "Pemasangan suspensi Coilover independen baru, gardan 9-inch Ford, dan restorasi rem cakram Wilwood di 4 roda.",
        image: "https://images.unsplash.com/photo-1486006920555-c77dce18193b?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t2-5",
        status: "Mesin V8 Big Block & Exhaust",
        date: "2026-07-10",
        description: "Instalasi mesin 427 V8 Shelby Stroker berkekuatan 500hp lengkap dengan transmisi manual Tremec 5-speed.",
        completed: false
      },
      {
        id: "t2-6",
        status: "Interior Re-upholstery & AC",
        date: "2026-08-05",
        description: "Instalasi roll cage Shelby, jok kulit hitam dengan sabuk pengaman 4 titik, serta penambahan unit AC Vintage Air.",
        completed: false
      },
      {
        id: "t2-7",
        status: "Final Road Test & Dyno Check",
        date: "2026-08-20",
        description: "Pengecekan output mesin di atas mesin Dyno untuk mengoptimalkan pembakaran karburator Holley, dan uji jalan sejauh 50 km.",
        completed: false
      }
    ],
    gallery: [
      {
        id: "g2-1",
        url: "https://images.unsplash.com/photo-1583121274602-3e2820c69888?auto=format&fit=crop&q=80&w=800",
        title: "Kondisi Awal Mobil",
        date: "2026-01-05",
        stage: "Intake",
        type: "photo"
      },
      {
        id: "g2-2",
        url: "https://images.unsplash.com/photo-1605558202076-168223c9f02b?auto=format&fit=crop&q=80&w=800",
        title: "Perbaikan Plat Lantai",
        date: "2026-02-28",
        stage: "Metalwork",
        type: "photo"
      },
      {
        id: "g2-video",
        url: "https://www.stanleymarthin.com/wp-content/uploads/2018/04/StanleyMarthin2.mp4",
        title: "Test Dyno & Suara Knalpot V8",
        date: "2026-03-15",
        stage: "Drivetrain",
        type: "video"
      },
      {
        id: "g2-3",
        url: "https://images.unsplash.com/photo-1615906655593-ad0386982a0f?auto=format&fit=crop&q=80&w=800",
        title: "Selesai Cat Eleanor Gray",
        date: "2026-04-12",
        stage: "Painting",
        type: "photo"
      },
      {
        id: "g2-4",
        url: "https://images.unsplash.com/photo-1605558202076-168223c9f02b?auto=format&fit=crop&q=80&w=800",
        title: "Instalasi Suspensi Wilwood",
        date: "2026-06-02",
        stage: "Suspension",
        type: "photo"
      }
    ]
  },
  {
    id: "v3",
    clientCode: "SM-ALEXA",
    clientName: "Alexandra Wijaya",
    brand: "Mercedes-Benz",
    model: "300 SL Gullwing",
    year: 1955,
    licensePlate: "W 1955 SL",
    vin: "1980405500050",
    status: "Completed & Delivered",
    completionPercentage: 100,
    estimatedCompletion: "2026-06-20",
    bannerImage: "https://images.unsplash.com/photo-1563720223185-11003d516935?auto=format&fit=crop&q=80&w=1200",
    restorationType: "Restorasi Total",
    timeline: [
      {
        id: "t3-1",
        status: "Registrasi & Dokumentasi Orisinil",
        date: "2025-06-15",
        description: "Verifikasi sertifikat Mercedes-Benz Classic Center Jerman untuk menjamin keaslian nomor rangka dan blok mesin.",
        image: "https://images.unsplash.com/photo-1563720223185-11003d516935?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t3-2",
        status: "Restorasi Rangka Tubular",
        date: "2025-09-10",
        description: "Rangka spaceframe tubular yang unik dibersihkan dari karat mikro dan dilapis powder coat abu-abu satin.",
        completed: true
      },
      {
        id: "t3-3",
        status: "Pengecatan Silver Metallic Orisinil",
        date: "2025-12-05",
        description: "Mengecat bodi aluminium-baja dengan cat Mercedes-Benz Silver Metallic (DB 180) yang legendaris.",
        completed: true
      },
      {
        id: "t3-4",
        status: "Instalasi Mesin M198 Fuel Injection",
        date: "2026-03-20",
        description: "Pemasangan mesin inline-six 3.0L injeksi langsung pertama di dunia yang telah dipoles dan dikalibrasi komponen internalnya.",
        completed: true
      },
      {
        id: "t3-5",
        status: "Upholstery Kulit Merah & Tartan",
        date: "2026-05-15",
        description: "Pengerjaan jok berlapis kulit merah berkualitas tertinggi dari Hans Reinke Jerman dan aksen interior orisinil.",
        completed: true
      },
      {
        id: "t3-6",
        status: "Uji Jalan & Serah Terima",
        date: "2026-06-20",
        description: "Selesai uji coba jalan sejauh 250 km. Mobil diserahkan kembali ke Ibu Alexandra Wijaya dalam kondisi sempurna (Concours Condition).",
        image: "https://images.unsplash.com/photo-1563720223185-11003d516935?auto=format&fit=crop&q=80&w=800",
        completed: true
      }
    ],
    gallery: [
      {
        id: "g3-1",
        url: "https://images.unsplash.com/photo-1563720223185-11003d516935?auto=format&fit=crop&q=80&w=800",
        title: "Tampak Samping Selesai Restorasi",
        date: "2026-06-20",
        stage: "Selesai",
        type: "photo"
      },
      {
        id: "g3-video",
        url: "https://www.stanleymarthin.com/wp-content/uploads/2018/04/StanleyMarthin2.mp4",
        title: "First Start Mesin Gullwing",
        date: "2026-03-20",
        stage: "Engine",
        type: "video"
      },
      {
        id: "g3-2",
        url: "https://images.unsplash.com/photo-1503376780353-7e6692767b70?auto=format&fit=crop&q=80&w=800",
        title: "Intake Awal di Bengkel",
        date: "2025-06-15",
        stage: "Intake",
        type: "photo"
      }
    ]
  },
  {
    id: "v4",
    clientCode: "SM-ADIT",
    clientName: "Varrel Aditya",
    brand: "Gowes VIP",
    model: "Becak Custom Sultan",
    year: 1998,
    licensePlate: "B 3 CAK",
    vin: "BCK-99887766",
    status: "Pemasangan Mesin V12 Twin Turbo Charge beserta Speaker JBL",
    completionPercentage: 30,
    estimatedCompletion: "2026-12-30",
    bannerImage: "/images/becak_2.jpg",
    restorationType: "Restorasi Sultan",
    timeline: [
      {
        id: "t4-1",
        status: "Inspeksi Masuk",
        date: "2026-06-01",
        description: "Becak tiba di bengkel dalam keadaan rantai putus dan jok sobek. Mang Udin sang pemilik menangis haru karena becaknya akan direstorasi menggunakan standar Porsche.",
        completed: true
      },
      {
        id: "t4-2",
        status: "Pemasangan Jok Kulit Nappa",
        date: "2026-06-15",
        description: "Jok penumpang diganti dengan kulit Nappa asli jahitan tangan agar penumpang tidak kepanasan.",
        image: "https://images.unsplash.com/photo-1486006920555-c77dce18193b?auto=format&fit=crop&q=80&w=800",
        completed: true
      },
      {
        id: "t4-3",
        status: "Pemasangan Speaker JBL",
        date: "2026-06-25",
        description: "Biar penumpang gak bosen di jalan, kita pasang speaker JBL lengkap dengan subwoofer di bawah jok penumpang.",
        completed: false
      }
    ],
    gallery: [
      {
        id: "g4-1",
        url: "/images/becak_2.jpg",
        title: "Kondisi Awal Becak",
        date: "2026-06-01",
        stage: "Intake",
        type: "photo"
      },
      {
        id: "g4-2",
        url: "/images/becak_1.jpg",
        title: "Referensi Warna Cat",
        date: "2026-06-15",
        stage: "Eksterior",
        type: "photo"
      },
      {
        id: "g4-3",
        url: "/images/becak_2.jpg",
        title: "Ide Modifikasi Motor",
        date: "2026-06-20",
        stage: "Mesin",
        type: "photo"
      },
      {
        id: "g4-4",
        url: "/images/becak_3.jpg",
        title: "Desain Kanopi Anti Hujan",
        date: "2026-06-25",
        stage: "Eksterior",
        type: "photo"
      },
      {
        id: "g4-5",
        url: "/images/becak_4.jpg",
        title: "Inspirasi Posisi Duduk",
        date: "2026-06-25",
        stage: "Interior",
        type: "photo"
      }
    ]
  },
  {
    id: "v5",
    clientCode: "SM-ADIT",
    clientName: "Aditya Pratama",
    brand: "TVS Racing",
    model: "Bajaj Roda Tiga Super",
    year: 2005,
    licensePlate: "B 4 JAJ",
    vin: "BJJ-11223344",
    status: "Ganti Knalpot Racing",
    completionPercentage: 80,
    estimatedCompletion: "2026-08-17",
    bannerImage: "/images/bajaj_4.jpg",
    restorationType: "Restorasi Kebut",
    timeline: [
      {
        id: "t5-1",
        status: "Bongkar Mesin 2-Tak",
        date: "2026-05-10",
        description: "Mesin 2-tak yang berisik dibongkar total. Suaranya bikin tetangga bengkel protes, jadi kita kalibrasi ulang agar suaranya sehalus V8.",
        completed: true
      },
      {
        id: "t5-2",
        status: "Pasang NOS & Turbo",
        date: "2026-06-25",
        description: "Biar bisa ngejar setoran dengan cepat, Bajaj ini dipasangin sistem NOS rahasia di bawah jok supir.",
        image: "https://images.unsplash.com/photo-1619642751034-765dfdf7c58e?auto=format&fit=crop&q=80&w=800",
        completed: true
      }
    ],
    gallery: [
      {
        id: "g5-1",
        url: "/images/bajaj_racing.png",
        title: "Bajaj Racing Siap Tempur",
        date: "2026-06-25",
        stage: "Modifikasi",
        type: "photo"
      },
      {
        id: "g5-2",
        url: "/images/bajaj_1.jpg",
        title: "Test Drive Kecepatan Tinggi",
        date: "2026-06-26",
        stage: "Test",
        type: "photo"
      },
      {
        id: "g5-3",
        url: "/images/bajaj_2.jpg",
        title: "Desain Body Kit",
        date: "2026-06-27",
        stage: "Eksterior",
        type: "photo"
      },
      {
        id: "g5-4",
        url: "/images/bajaj_3.jpg",
        title: "Inspirasi Cat Ungu Premium",
        date: "2026-06-28",
        stage: "Paint",
        type: "photo"
      },
      {
        id: "g5-5",
        url: "/images/bajaj_4.jpg",
        title: "Pengecekan Kaki-kaki",
        date: "2026-06-29",
        stage: "Mesin",
        type: "photo"
      }
    ]
  }
];
