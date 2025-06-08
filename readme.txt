Usage:
  txtr <input> <output> [options]

Options:
  --kvinput:  Input for a keyvalue file

              kvinput takes '=' separated values and substitutes them into ${<key>} from the input file, then outputs

Examples:
  txtr docker-compose.example.yml docker-compose.yml --kvinput .env