#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[POLIKLINIK GIGI & MULUT] \
      #text(10pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Rekam Medis Odontogram - Terverifikasi Dokter Gigi]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9pt)

#align(center)[
  #text(12pt, weight: "bold")[REKAM MEDIS KESEHATAN GIGI & ODONTOGRAM] \
  #text(8.5pt)[Nomor Dokumen: {{DOC_NUMBER}}]
]

#v(0.5em)

#table(
  columns: (120pt, 8pt, 1fr, 100pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Pemeriksaan], [:], [{{ISSUED_DATE}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Dokter Gigi Pemeriksa], [:], [{{DOCTOR_NAME}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Keluhan Utama], [:], [{{KELUHAN_UTAMA}}],
)

#v(0.5em)
*PEMETAAN ODONTOGRAM (STANDAR FDI DENTAL NOTATION):*
#table(
  columns: (1fr, 1fr),
  stroke: 0.5pt,
  [
    #align(center)[*Kuadran 1 (Kanan Atas)*] \
    18 - 17 - 16 - 15 - 14 - 13 - 12 - 11 \
    Status: {{KUADRAN_1}}
  ],
  [
    #align(center)[*Kuadran 2 (Kiri Atas)*] \
    21 - 22 - 23 - 24 - 25 - 26 - 27 - 28 \
    Status: {{KUADRAN_2}}
  ],
  [
    #align(center)[*Kuadran 4 (Kanan Bawah)*] \
    48 - 47 - 46 - 45 - 44 - 43 - 42 - 41 \
    Status: {{KUADRAN_4}}
  ],
  [
    #align(center)[*Kuadran 3 (Kiri Bawah)*] \
    31 - 32 - 33 - 34 - 35 - 36 - 37 - 38 \
    Status: {{KUADRAN_3}}
  ],
)

#v(0.5em)
*PEMERIKSAAN INTRAORAL LAINNYA:*
#table(
  columns: (1fr, 1fr, 1fr),
  stroke: 0.5pt,
  [ *Oklusi:* {{OKLUSI}} ],
  [ *Torus Palatinus:* {{TORUS_PALATINUS}} ],
  [ *Torus Mandibularis:* {{TORUS_MANDIBULARIS}} ],
  [ *Palatum:* {{PALATUM}} ],
  [ *Diastema:* {{DIASTEMA}} ],
  [ *Gigi Anomali:* {{ANOMALI_GIGI}} ],
)

#v(0.5em)
*DIAGNOSA GIGI & RENCANA TINDAKAN (TREATMENT PLAN):*
#table(
  columns: (80pt, 140pt, 1fr),
  stroke: 0.5pt,
  table.header([*Elemen Gigi*], [*Diagnosa Klinis*], [*Rencana Tindakan (Scaling, Tumpatan, Exo)*]),
  {{ODONTOGRAM_TREATMENT_ROWS}}
)

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI DIGITAL DOKTER GIGI] \
        #text(6.5pt)[Pindai kode QR untuk validasi odontogram:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Gigi Pemeriksa, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
