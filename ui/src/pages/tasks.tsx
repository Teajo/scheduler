import 'date-fns';
import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableContainer from '@material-ui/core/TableContainer';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';
import Picker from './pickers';
import axios from 'axios';

interface Column {
  id: string;
  label: string;
  minWidth?: number;
  align?: 'right';
  format?: (value: number) => string;
}

const columns: Column[] = [
  { id: 'id', label: 'ID', minWidth: 170 },
  { id: 'date', label: 'Date', minWidth: 100 },
  { id: 'done', label: 'Done', minWidth: 100, format: l => String(l) },
];

const useStyles = makeStyles({
  root: {
    width: '100%',
  },
  container: {
    maxHeight: 700,
  },
});

export default function StickyHeadTable(props: any) {
  const classes = useStyles();
  const [page, setPage] = React.useState(0);
  const [rowsPerPage, setRowsPerPage] = React.useState(10);
  const [selectedStartDate, setSelectedStartDate] = React.useState<Date>(new Date());
  const [selectedEndDate, setSelectedEndDate] = React.useState<Date>(new Date());
  const [tasks, setTasks] = React.useState([]);

  React.useEffect(() => {
    axios.get(`http://127.0.0.1:3000/tasks?endDate=${selectedEndDate.toISOString()}&startDate=${selectedStartDate.toISOString()}`)
    .then(res => {
      const { data } = res.data;
      setTasks(data);
    })
    .catch(err => {
      console.log(err);
    })
  }, [selectedStartDate, selectedEndDate]);

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  const onStartChanged = (date: Date | null) => {
    if (date) {
      setSelectedStartDate(date);
    }
  };

  const onEndChanged = (date: Date | null) => {
    if (date) {
      setSelectedEndDate(date);
    }
  };

  return (
    <>
      <div style={{ textAlign: 'left' }}>
        <Picker label="Start" date={selectedStartDate} onChange={onStartChanged} />&nbsp;&nbsp;&nbsp;
        <Picker label="End" date={selectedEndDate} onChange={onEndChanged} />
      </div>

      <br />
      
      <Paper className={classes.root}>
        <TableContainer className={classes.container}>
          <Table stickyHeader aria-label="sticky table">
            <TableHead>
              <TableRow>
                {columns.map((column) => (
                  <TableCell
                    key={column.id}
                    align={column.align}
                    style={{ minWidth: column.minWidth }}
                  >
                    {column.label}
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {tasks.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage).map((row: any) => {
                return (
                  <TableRow hover role="checkbox" tabIndex={-1} key={row.code}>
                    {columns.map((column) => {
                      const value = row[column.id];
                      return (
                        <TableCell key={column.id} align={column.align}>
                          {column.format ? column.format(value) : value}
                        </TableCell>
                      );
                    })}
                  </TableRow>
                );
              })}
            </TableBody>
          </Table>
        </TableContainer>
        <TablePagination
          rowsPerPageOptions={[10, 25, 100]}
          component="div"
          count={tasks.length}
          rowsPerPage={rowsPerPage}
          page={page}
          onChangePage={handleChangePage}
          onChangeRowsPerPage={handleChangeRowsPerPage}
        />
      </Paper>
    </>
  );
}
