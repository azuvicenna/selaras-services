#set page(
  paper: "a4",
  margin: (x: 2cm, y: 2.5cm),
  header: [
    #align(center)[
      #text(13pt, weight: "bold")[RUMAH SAKIT SELARAS MEDIKA] \
      #text(8pt)[Jl. Kesehatan Sejahtera No. 100, Jakarta | Telp: (021) 555-0123 | Email: rs\@selaras.hospital.id]
      #v(2pt)
      #line(length: 100%, stroke: 1.5pt)
    ]
  ],
  footer: [
    #line(length: 100%, stroke: 0.5pt)
    #grid(
      columns: (1fr, 1fr),
      align(left)[#text(7pt, fill: luma(100))[Dokumen Resmi RS Selaras Medika - Terverifikasi Elektronik]],
      align(right)[#text(7pt)[Halaman #context counter(page).display()]]
    )
  ]
)

#set text(size: 10pt)

#align(center)[
  #text(12pt, weight: "bold")[SURAT KETERANGAN RAWAT INAP] \
  #text(9pt)[Nomor: {{DOC_NUMBER}}]
]

#v(1em)

Yang bertanda tangan di bawah ini dokter Rumah Sakit Selaras Medika menerangkan bahwa:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: none,
  [Nama Pasien], [:], [{{PATIENT_NAME}}],
  [Nomor Rekam Medis], [:], [{{PATIENT_ID}}],
  [Jenis Kelamin / Tgl Lahir], [:], [{{PATIENT_GENDER}} / {{PATIENT_DOB}}],
  [Alamat], [:], [{{ALAMAT}}],
)

#v(0.5em)
Adalah benar pasien tersebut di atas telah menjalani / sedang menjalani perawatan rawat inap di Rumah Sakit Selaras Medika dengan rincian:

#table(
  columns: (140pt, 10pt, 1fr),
  stroke: 0.5pt,
  [Tanggal Masuk Rawat], [:], [{{TGL_MASUK}}],
  [Tanggal Keluar Rawat], [:], [{{TGL_KELUAR}}],
  [Ruang / Kelas Perawatan], [:], [{{RUANG_RAWAT}} / {{KELAS_RAWAT}}],
  [Diagnosa Medis], [:], [{{DIAGNOSIS}}],
  [Dokter Penanggung Jawab], [:], [{{DOCTOR_NAME}}],
)

#v(1em)
Surat keterangan ini diberikan atas permintaan pasien / keluarga untuk keperluan: *{{KEPERLUAN}}*.

#v(1.5em)
Demikian surat keterangan ini dibuat dengan sebenarnya untuk dipergunakan sebagaimana mestinya.

#v(2em)

#grid(
  columns: (1fr, 1fr),
  [
    #rect(stroke: 0.5pt, inset: 6pt)[
      #align(center)[
        #text(7pt, weight: "bold")[VERIFIKASI KEASLIAN ELEKTRONIK] \
        #text(6.5pt)[Pindai kode QR atau akses URL di bawah untuk verifikasi keabsahan dokumen:] \
        #v(2pt)
        #text(6.5pt, font: ("Courier New", "DejaVu Sans Mono", "monospace"))[{{VERIFY_URL}}]
      ]
    ]
  ],
  [
    #align(center)[
      Jakarta, {{ISSUED_DATE}} \
      Dokter Pemeriksa / DPJP, \
      #v(3.5em)
      #text(weight: "bold")[{{DOCTOR_NAME}}] \
      SIP: {{DOCTOR_SIP}}
    ]
  ]
)
