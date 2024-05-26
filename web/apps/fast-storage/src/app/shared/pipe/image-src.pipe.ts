import { Pipe, PipeTransform } from '@angular/core';

type FileTypeToImageMap = {
  [key: string]: string;
};

@Pipe({
  name: 'imageSrc',
  standalone: true,
})
export class ImageSrcPipe implements PipeTransform {
  transform(path: string, args: any[]): string {
    const fileTypeToImage: FileTypeToImageMap = {
      PDF: 'pdf.png',
      DOCX: 'docx.png',
      JPG: 'jpg.png',
      CSV: 'misc.png',
      PNG: 'jpg.png',
      EXE: 'exe.png',
      ZIP: 'zip.png',
      RAR: 'zip.png',
      MSI: 'misc.png',
      XLSX: 'xlsx.png',
      MSG: 'misc.png',
      JS: 'js.png',
    };

    const [fileType, isDir] = args;
    const defaultImage = isDir ? 'folder.png' : 'misc.png';
    const imagePath = fileTypeToImage[fileType] || defaultImage;
    return `${path}${imagePath}`;
  }
}
