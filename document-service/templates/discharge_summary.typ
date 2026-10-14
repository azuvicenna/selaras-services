#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.2cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Ringkasan Pulang Pasien (Discharge Summary) - Rahasia Medis]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 9.5pt)

#align(center)[
  #text(12pt, weight: "bold")[RINGKASAN PULANG PASIEN (DISCHARGE SUMMARY)] \
  #text(9pt)[Nomor Dokumen: {{DOC_NUMBER}}]
]

#v(0.5em)

#table(
  columns: (120pt, 8pt, 1fr, 100pt, 8pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}], [Tgl Masuk], [:], [{{TGL_MASUK}}],
  [No. Rekam Medis], [:], [{{PATIENT_ID}}], [Tgl Keluar], [:], [{{TGL_KELUAR}}],
  [Jenis Kelamin / Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}], [Ruang Rawat], [:], [{{RUANG_RAWAT}}],
  [Dokter DPJP], [:], [{{DOCTOR_NAME}}], [Cara Keluar], [:], [{{CARA_KELUAR}}],
)

#v(0.5em)

#table(
  columns: (1fr),
  stroke: 0.5pt,
  [ *1. Diagnosa Masuk Rawat Inap:* \ {{DIAGNOSIS_MASUK}} ],
  [ *2. Diagnosa Utama (Final) & Sekunder (ICD-10):* \ {{DIAGNOSIS_UTAMA}} \ {{DIAGNOSIS_SEKUNDER}} ],
  [ *3. Ringkasan Riwayat Penyakit & Pemeriksaan Fisik:* \ {{RINGKASAN_KLINIS}} ],
  [ *4. Hasil Pemeriksaan Penunjang Penting (Lab, Radiologi, EKG):* \ {{PEMERIKSAAN_PENUNJANG}} ],
  [ *5. Tindakan Medis / Operasi / Prosedur (ICD-9-CM):* \ {{PROSEDUR}} ],
  [ *6. Kondisi Saat Pulang:* \ {{KONDISI_PULANG}} ],
  [ *7. Terapi Pulang / Obat Lanjutan:* \ {{OBAT_PULANG}} ],
  [ *8. Instruksi Lanjutan / Tanggal Kontrol Ulang:* \ {{KONTROL_ULANG}} ],
)

#v(1.5em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI DIGITAL DPJP] \
        #text(6.5pt)[Telah ditinjau dan disetujui secara elektronik:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Penanggung Jawab Pelayanan (DPJP), \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
