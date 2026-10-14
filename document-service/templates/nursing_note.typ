#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(10pt, weight: "bold")[CATATAN PERKEMBANGAN PASIEN TERINTEGRASI (CPPT)] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Catatan Asuhan Keperawatan Resmi - RS Selaras Medika]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9pt)

#align(center)[
  #text(12pt, weight: "bold")[LEMBAR CATATAN OBSERVASI & ASUHAN KEPERAWATAN] \
  #text(8.5pt)[Nomor Dokumen: {{DOC_NUMBER}}]
]

#v(0.5em)

#table(
  columns: (120pt, 8pt, 1fr, 100pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tanggal Catatan], [:], [{{ISSUED_DATE}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Ruang Perawatan], [:], [{{RUANG_RAWAT}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Perawat Jaga], [:], [{{NURSE_NAME}}],
)

#v(0.5em)
*TANDA-TANDA VITAL & EARLY WARNING SCORE (EWS):*
#table(
  columns: (1fr, 1fr, 1fr, 1fr, 1fr),
  stroke: 0.5pt,
  [ *Tekanan Darah:* {{TD}} mmHg ],
  [ *Nadi:* {{NADI}} x/m ],
  [ *Suhu:* {{SUHU}} °C ],
  [ *Laju Napas:* {{PERNAPASAN}} x/m ],
  [ *SpO2 / Skala Nyeri:* {{SPO2}}% / {{NYERI}} ],
)

#v(0.5em)
*DOKUMENTASI ASUHAN KEPERAWATAN (SOAP):*
#table(
  columns: (90pt, 8pt, 1fr),
  stroke: 0.5pt,
  [S (Subjective)], [:], [{{SUBJEKTIF}}],
  [O (Objective)], [:], [{{OBJEKTIF}}],
  [A (Assessment)], [:], [{{MASALAH_KEPERAWATAN}}],
  [P (Planning)], [:], [{{INTERVENSI}}],
)

#v(0.5em)
*TINDAKAN / IMPLEMENTASI KEPERAWATAN YANG DILAKUKAN:*
#block(inset: (left: 10pt))[
  {{IMPLEMENTASI}}
]

#v(0.5em)
*EVALUASI HASIL & INSTRUKSI SHIFT BERIKUTNYA:*
#block(inset: (left: 10pt))[
  {{EVALUASI}}
]

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI SISTEM ELEKTRONIK] \
        #text(6.5pt)[Tervalidasi dalam basis data CPPT terintegrasi:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Perawat Pelaksana Asuhan, \
      #v(3.5em)
      #text(weight: "bold")[{{NURSE_NAME}}] \
      NIRA / STR: {{NURSE_STR}}
    ]
  ]
)
